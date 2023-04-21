package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ixpectus/declarate/suite"
	"gopkg.in/yaml.v2"
)

type Converter struct {
	sourceDir string
	targetDir string
}

func New(sourceDir, targetDir string) *Converter {
	return &Converter{
		sourceDir: sourceDir,
		targetDir: targetDir,
	}
}

func (c *Converter) Convert() error {
	s := suite.New(c.sourceDir, suite.RunConfig{})
	tt, err := s.AllTests(c.sourceDir)
	if err != nil {
		return err
	}
	for _, v := range tt {
		tests := []GonkeyTest{}
		relatedName := strings.ReplaceAll(v, c.sourceDir, "")
		if relatedName == "" {
			relatedName = path.Base(c.sourceDir)
		}
		data, err := os.ReadFile(v)
		if err != nil {
			log.Printf("failed to convert %s, %v", v, err)
			continue
		}
		err = yaml.Unmarshal(data, &tests)
		if err != nil {
			log.Printf("failed to unmarshall on convert %s, %v", v, err)
			continue
		}
		converted := convert(tests)
		if len(converted) == 0 {
			continue
		}
		bb, _ := yaml.Marshal(converted)
		targetFile := c.targetDir + "/" + relatedName

		if err := os.MkdirAll(path.Dir(targetFile), os.ModePerm); err != nil {
			log.Printf("failed to mkdir for file %s, %v", targetFile, err)
		}
		err = os.WriteFile(targetFile, bb, os.ModePerm)
		if err != nil {
			log.Printf("failed to write to file %s, %v", targetFile, err)
			continue
		}
	}

	return nil
}

func convert(originalTests []GonkeyTest) []DeclarateTest {
	res := make([]DeclarateTest, 0, len(originalTests))
	i := 0
	for _, v := range originalTests {
		i++
		converted := DeclarateTest{}
		converted.Name = v.Name
		if v.RequestURL != "" && (v.DbQueryTmpl != "" || len(v.DatabaseChecks) > 0 || v.AfterRequestScriptParams != nil) {
			convertedRequest := request(v)
			convertedRequest.Name += " in api"
			convertedDb := db(v)
			for k := range convertedDb {
				convertedDb[k].Name += " in db"
			}
			steps := []DeclarateTest{convertedRequest}
			steps = append(steps, convertedDb...)
			converted.Steps = steps
		} else if v.DbQueryTmpl != "" || len(v.DatabaseChecks) > 0 {
			if len(v.DatabaseChecks) > 0 {
				steps := []DeclarateTest{}
				convertedDb := db(v)
				steps = append(steps, convertedDb...)
				converted.Steps = steps
			} else if v.DbQueryTmpl != "" {
				converted = db(v)[0]
				converted.Name = v.Name
			}
		} else if v.RequestURL != "" {
			converted = request(v)
		} else if v.Variables != nil || len(v.Tags) > 0 {
			if len(v.Tags) > 0 {
				converted.Definition = &Definition{
					Tags: v.Tags,
				}
				res = append(res, converted)
			}
			if v.Variables != nil {
				converted = DeclarateTest{}
				converted.Name = v.Name
				converted = variables(v, converted)
				res = append(res, converted)
			}
			continue
		} else if len(v.Tags) > 0 {
			converted.Definition = &Definition{
				Tags: v.Tags,
			}
		} else if v.AfterRequestScriptParams != nil {
			converted = afterRequest(v)
		} else {
			continue
		}
		if v.Poll != nil {
			converted.Poll = &Poll{
				Duration: v.Poll.Duration,
			}
			if v.Poll.ResponseBodyRegexp != "" {
				converted.Poll.ResponseBodyRegexp = v.Poll.ResponseBodyRegexp
			}
			if v.Poll.Interval > 0 {
				converted.Poll.Interval = v.Poll.Interval
			}
		}
		if len(v.PollInterval) > 0 {
			duration := time.Duration(0)
			for _, v := range v.PollInterval {
				duration = duration + v
			}
			converted.Poll = &Poll{
				Duration: duration,
			}
		}
		res = append(res, converted)
	}
	return res
}

func afterRequest(g GonkeyTest) DeclarateTest {
	res := DeclarateTest{}
	res.Name = g.Name
	res.ScriptPath = g.AfterRequestScriptParams.PathTmpl
	return res
}

func request(g GonkeyTest) DeclarateTest {
	res := DeclarateTest{}
	res.Name = g.Name
	res.RequestTmpl = varFix(g.RequestTmpl)
	res.RequestURL = varFix(g.RequestURL)
	res.ComparisonParams = g.ComparisonParams
	res.Method = g.Method
	if g.HeadersVal != nil {
		res.HeadersVal = map[string]string{}
		for k, v := range g.HeadersVal {
			res.HeadersVal[k] = varFix(v)
		}
	}
	for _, v := range g.VariablesToSet {
		res.VariablesToSet = map[string]string{}
		for k1, v1 := range v {
			res.VariablesToSet[k1] = "body." + v1
		}
	}
	for k, v := range g.ResponseTmpls {
		js1 := fmt.Sprintf("{\"body\":%v,\"status\":%v}", v, k)
		var prettyJSON bytes.Buffer
		_ = json.Indent(&prettyJSON, []byte(js1), "", "  ")
		res.ResponseTmpls = prettyJSON.String()
		break
	}
	return res
}

func varFix(v string) string {
	v = strings.ReplaceAll(v, "{{ ", "{{")
	v = strings.ReplaceAll(v, " }}", "}}")
	return v
}

func db(g GonkeyTest) []DeclarateTest {
	tests := []DeclarateTest{}
	dbConn := varFix(g.DBConnectionString)
	if g.DbQueryTmpl != "" {
		d := DeclarateTest{
			Name:           g.Name,
			DbConn:         dbConn,
			DbQuery:        g.DbQueryTmpl,
			VariablesToSet: g.DBVariablesToSet,
		}
		if len(g.DbResponseTmpl) > 0 {
			d.DbResponse = "[" + strings.Join(g.DbResponseTmpl, ",") + "]"
		}
		tests = append(tests, d)
	}
	for i, v := range g.DatabaseChecks {
		d := DeclarateTest{
			Name:           g.Name + fmt.Sprintf("#%v", i),
			DbConn:         dbConn,
			DbQuery:        v.DbQueryTmpl,
			VariablesToSet: v.DBVariablesToSet,
		}
		if len(v.DbResponseTmpl) > 0 {
			d.DbResponse = "[" + strings.Join(v.DbResponseTmpl, ",") + "]"
		}
		tests = append(tests, d)
	}
	return tests
}

func variables(g GonkeyTest, res DeclarateTest) DeclarateTest {
	res.Name = g.Name
	res.Variables = map[string]string{}
	for k, v := range g.Variables {
		res.Variables[k] = strings.ReplaceAll(v, "$eval", "$")
	}
	return res
}
