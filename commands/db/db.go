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

type DbCommand struct {
	Config       *DbCheck
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

	if cfg == nil {
		return nil, nil
	}
	return &DbCommand{
		Config:     cfg,
		Connection: u.connection,
	}, nil
}

type DbCheck struct {
	Check *CheckConfig `yaml:"db,omitempty"`
}

type CheckConfig struct {
	DbConn           string                `json:"dbConn" yaml:"dbConn"`
	DbQuery          string                `json:"dbQuery" yaml:"dbQuery"`
	DbResponse       string                `json:"dbResponse" yaml:"dbResponse"`
	ComparisonParams compare.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
	VariablesToSet   map[string]string     `yaml:"variables_to_set"`
}

func (e *DbCommand) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *DbCommand) Do() error {
	if e.Config.Check != nil {
		db, err := sql.Open("postgres", e.Connection)
		if err != nil {
			return err
		}

		if e.Config.Check.DbResponse != "" || e.Config.Check.VariablesToSet != nil {
			res, err := makeQuery(e.Config.Check.DbQuery, db)
			if err != nil {
				return err
			}
			e.responseBody = &res

			return nil
		}
		if err := execQuery(e.Config.Check.DbQuery, db); err != nil {
			return err
		}
	}

	return nil
}

func (e *DbCommand) ResponseBody() *string {
	return e.responseBody
}

func (e *DbCommand) VariablesToSet() map[string]string {
	if e != nil && e.Config != nil {
		return e.Config.Check.VariablesToSet
	}
	return nil
}

func (e *DbCommand) Check() error {
	if e.Config.Check != nil && e.responseBody != nil && e.Config.Check.DbResponse != "" {
		errs, err := compareJsonBody(e.Config.Check.DbResponse, *e.responseBody, e.Config.Check.ComparisonParams)
		if len(errs) > 0 {
			msg := ""
			for _, v := range errs {
				msg += v.Error() + "\n"
			}
			return &contract.TestError{
				Title:         "response body differs",
				Expected:      e.Config.Check.DbResponse,
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
