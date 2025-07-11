package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dailymotion/allure-go"
	_ "github.com/lib/pq"
	"github.com/xwb1989/sqlparser"

	"github.com/ixpectus/declarate/contract"
)

type Db struct {
	Config        *CheckConfig
	Vars          contract.Vars
	Report        contract.ReportAttachement
	Comparer      contract.Comparer
	connectLoader contract.DBConnectLoader
	responseBody  *string
}

type Unmarshaller struct {
	connectLoader contract.DBConnectLoader
	comparer      contract.Comparer
}

func NewUnmarshaller(
	connectLoader contract.DBConnectLoader,
	comparer contract.Comparer,
) *Unmarshaller {
	return &Unmarshaller{
		connectLoader: connectLoader,
		comparer:      comparer,
	}
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
			Config:        cfg.Check,
			connectLoader: u.connectLoader,
			Comparer:      u.comparer,
		}, nil
	}

	return &Db{
		Config:        cfgShort,
		connectLoader: u.connectLoader,
		Comparer:      u.comparer,
	}, nil
}

type DbCheck struct {
	Check *CheckConfig `yaml:"db,omitempty"`
}

func (d *DbCheck) isEmpty() bool {
	return d == nil || d.Check == nil || d.Check.DbQuery == ""
}

type CheckConfig struct {
	DbConn           string                 `json:"db_conn" yaml:"db_conn"`
	DbQuery          string                 `json:"db_query" yaml:"db_query"`
	DbResponse       string                 `json:"db_response" yaml:"db_response"`
	ComparisonParams contract.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
}

func (d *CheckConfig) isEmpty() bool {
	return d == nil || d.DbQuery == ""
}

func (e *Db) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *Db) SetReport(r contract.ReportAttachement) {
	e.Report = r
}

func (e *Db) GetConfig() interface{} {
	return e.Config
}

func (e *Db) Do() error {
	if e.Config != nil {
		e.Config.DbConn = e.Vars.Apply(e.Config.DbConn)

		e.Config.DbQuery = e.Vars.Apply(e.Config.DbQuery)

		if e.Report != nil {
			e.Report.AddAttachment("query", allure.TextPlain, []byte(e.Config.DbQuery))
			if e.Config.DbConn != "" {
				e.Report.AddAttachment("conn", allure.TextPlain, []byte(e.Config.DbConn))
			} else {
				e.Report.AddAttachment("conn", allure.TextPlain, []byte(e.connectLoader.DefaultConnectionString()))
			}
		}
		e.Config.DbResponse = e.Vars.Apply(e.Config.DbResponse)
		db, err := e.connectLoader.Get(e.Config.DbConn)
		if err != nil {
			return err
		}
		defer db.Close()
		isSelect, err := isSelectStatement(e.Config.DbQuery)
		if err != nil {
			return err
		}

		if isSelect {
			res, err := makeQuery(e.Config.DbQuery, db)
			if err != nil {
				return err
			}
			e.responseBody = &res
			e.Report.AddAttachment("response", allure.TextPlain, []byte(res))
			return nil
		}
		if err := execQuery(e.Config.DbQuery, db); err != nil {
			return err
		}
	}

	return nil
}

func (e *Db) IsValid() error {
	if e.Config.DbResponse != "" {
		valid := json.Valid([]byte(e.Config.DbResponse))
		if !valid {
			return fmt.Errorf("cannot parse db response: `%v`", e.Config.DbResponse)
		}
	}
	return nil
}

func (e *Db) ResponseBody() *string {
	return e.responseBody
}

func (e *Db) Check() error {
	if e.Config != nil && e.responseBody != nil && e.Config.DbResponse != "" {
		errs, err := e.Comparer.CompareJsonBody(
			e.Config.DbResponse,
			*e.responseBody,
			e.Config.ComparisonParams,
		)
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

func execQuery(dbQuery string, db *sql.DB) error {
	queries := strings.Split(dbQuery, ";")
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return fmt.Errorf("failed exec db query %s, err %v", q, err)
		}
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

func isSelectStatement(dbQuery string) (bool, error) {
	queries := strings.Split(dbQuery, ";")
	stmt, err := sqlparser.Parse(queries[0])
	if err != nil {
		dbQuery = strings.Trim(dbQuery, " ")
		dbQuery = strings.ToLower(dbQuery)
		return strings.HasPrefix(dbQuery, "select"), nil
	}
	switch stmt.(type) {
	case *sqlparser.Select:
		return true, nil
	case *sqlparser.Union:
		return true, nil
	case *sqlparser.ParenSelect:
		return true, nil
	default:
		return false, nil
	}
}
