package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"

	"github.com/kylelemons/godebug/pretty"
)

type Db struct {
	Config       *CheckConfig
	Vars         contract.Vars
	Connection   string
	responseBody *string
}

type Unmarshaller struct {
	connection string
}

func NewUnmarshaller(connection string) *Unmarshaller {
	return &Unmarshaller{connection: connection}
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfg := &DbCheck{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}
	cfgShort := &CheckConfig{}
	if err := unmarshal(cfgShort); err != nil {
		return nil, err
	}
	if cfg.isEmpty() && cfgShort.isEmpty() {
		return nil, nil
	}
	if !cfg.Check.isEmpty() {
		return &Db{
			Config:     cfg.Check,
			Connection: u.connection,
		}, nil
	}
	return &Db{
		Config:     cfgShort,
		Connection: u.connection,
	}, nil
}

type DbCheck struct {
	Check *CheckConfig `yaml:"db,omitempty"`
}

func (d *DbCheck) isEmpty() bool {
	return d == nil || d.Check == nil || d.Check.DbQuery == ""
}

type CheckConfig struct {
	DbConn           string                `json:"dbConn" yaml:"dbConn"`
	DbQuery          string                `json:"dbQuery" yaml:"dbQuery"`
	DbResponse       string                `json:"dbResponse" yaml:"dbResponse"`
	ComparisonParams compare.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
	VariablesToSet   map[string]string     `yaml:"variables_to_set"`
}

func (d *CheckConfig) isEmpty() bool {
	return d == nil || d.DbQuery == ""
}

func (e *Db) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *Db) Do() error {
	if e.Config != nil {
		e.Config.DbConn = e.Vars.Apply(e.Config.DbConn)
		e.Config.DbQuery = e.Vars.Apply(e.Config.DbQuery)
		e.Config.DbResponse = e.Vars.Apply(e.Config.DbResponse)
		db, err := sql.Open("postgres", e.Connection)
		if err != nil {
			return err
		}

		if e.Config.DbResponse != "" || e.Config.VariablesToSet != nil {
			res, err := makeQuery(e.Config.DbQuery, db)
			if err != nil {
				return err
			}
			e.responseBody = &res

			return nil
		}
		if err := execQuery(e.Config.DbQuery, db); err != nil {
			return err
		}
	}

	return nil
}

func (e *Db) ResponseBody() *string {
	return e.responseBody
}

func (e *Db) VariablesToSet() map[string]string {
	if e != nil && e.Config != nil {
		return e.Config.VariablesToSet
	}
	return nil
}

func (e *Db) Check() error {
	if e.Config != nil && e.responseBody != nil && e.Config.DbResponse != "" {
		errs, err := compare.CompareJsonBody(e.Config.DbResponse, *e.responseBody, e.Config.ComparisonParams)
		if len(errs) > 0 {
			msg := ""
			for _, v := range errs {
				msg += v.Error() + "\n"
			}
			return &contract.TestError{
				Title:         "response body differs",
				Expected:      e.Config.DbResponse,
				Actual:        *e.responseBody,
				Message:       msg,
				OriginalError: fmt.Errorf("response body differs: %v", msg),
			}
		}
		if err != nil {
			return fmt.Errorf("compare json failed: %w", err)
		}
	}

	return nil
}

func compareJsonBody(expectedBody string, realBody string, params compare.CompareParams) ([]error, error) {
	// decode expected body
	var expected interface{}
	if err := json.Unmarshal([]byte(expectedBody), &expected); err != nil {
		return nil, fmt.Errorf(
			"invalid JSON in response for test : %s",
			err.Error(),
		)
	}

	// decode actual body
	var actual interface{}
	if err := json.Unmarshal([]byte(realBody), &actual); err != nil {
		return []error{errors.New("could not parse response")}, nil
	}

	return compare.Compare(expected, actual, params), nil
}

func toJsonArray(items []string, qual, testName string) ([]interface{}, error) {
	var itemJSONs []interface{}
	for i, row := range items {
		var itemJson interface{}
		if err := json.Unmarshal([]byte(row), &itemJson); err != nil {
			return nil, fmt.Errorf(
				"invalid JSON in the %s DB response for test %s:\n row #%d:\n %s\n error:\n%s",
				qual,
				testName,
				i,
				row,
				err.Error(),
			)
		}
		itemJSONs = append(itemJSONs, itemJson)
	}
	return itemJSONs, nil
}

func compareDbResponseLength(expected, actual []string, query interface{}) error {
	var err error

	if len(expected) != len(actual) {
		err = fmt.Errorf(
			"quantity of items in database do not match (-expected: %s +actual: %s)\n     test query:\n%s\n    result diff:\n%s",
			color.CyanString("%v", len(expected)),
			color.CyanString("%v", len(actual)),
			color.CyanString("%v", query),
			color.CyanString("%v", pretty.Compare(expected, actual)),
		)
	}
	return err
}

func execQuery(dbQuery string, db *sql.DB) error {
	if idx := strings.IndexByte(dbQuery, ';'); idx >= 0 {
		dbQuery = dbQuery[:idx]
	}
	if _, err := db.Exec(dbQuery); err != nil {
		return err
	}

	return nil
}

func makeQuery(dbQuery string, db *sql.DB) (string, error) {
	var dbResponse []string
	var jsonString string

	if idx := strings.IndexByte(dbQuery, ';'); idx >= 0 {
		dbQuery = dbQuery[:idx]
	}

	rows, err := db.Query(fmt.Sprintf("SELECT row_to_json(rows) FROM (%s) rows;", dbQuery))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&jsonString)
		if err != nil {
			return "", err
		}
		dbResponse = append(dbResponse, jsonString)
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	result := "[" + strings.Join(dbResponse, ", ") + "]"

	return result, nil
}