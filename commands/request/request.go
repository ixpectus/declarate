package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
)

type Request struct {
	Config         *RequestConfig
	Vars           contract.Vars
	Host           string
	responseBody   *string
	responseStatus int
}

type Unmarshaller struct {
	host string
}

func NewUnmarshaller(host string) *Unmarshaller {
	return &Unmarshaller{host: host}
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfg := &RequestConfig{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, nil
	}
	return &Request{
		Config: cfg,
		Host:   u.host,
	}, nil
}

type RequestConfig struct {
	Method           string                    `json:"method" yaml:"method"`
	RequestTmpl      string                    `json:"request" yaml:"request"`
	ResponseTmpls    string                    `json:"response" yaml:"response"`
	ResponseStatus   *int                      `json:"responseStatus" yaml:"responseStatus"`
	ResponseHeaders  map[int]map[string]string `json:"responseHeaders" yaml:"responseHeaders"`
	HeadersVal       map[string]string         `json:"headers" yaml:"headers"`
	QueryParams      string                    `json:"query" yaml:"query"`
	CookiesVal       map[string]string         `json:"cookies" yaml:"cookies"`
	ComparisonParams compare.CompareParams     `json:"comparisonParams" yaml:"comparisonParams"`
	RequestURL       string                    `json:"path" yaml:"path"`
	VariablesToSet   map[string]string         `yaml:"variables_to_set"`
}

func (e *Request) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *Request) Do() error {
	if e.Config.Method != "" {
		e.Config.QueryParams = e.Vars.Apply(e.Config.QueryParams)
		e.Config.RequestTmpl = e.Vars.Apply(e.Config.RequestTmpl)
		e.Config.ResponseTmpls = e.Vars.Apply(e.Config.ResponseTmpls)
		e.Config.RequestURL = e.Vars.Apply(e.Config.RequestURL)
		req, err := newCommonRequest(e.Host, *e.Config)
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		e.responseStatus = resp.StatusCode
		_ = resp.Body.Close()
		if err != nil {
			return err
		}
		s := string(body)
		e.responseBody = &s
	}

	return nil
}

func (e *Request) ResponseBody() *string {
	return e.responseBody
}

func (e *Request) VariablesToSet() map[string]string {
	if e != nil && e.Config != nil {
		return e.Config.VariablesToSet
	}
	return nil
}

func (e *Request) Check() error {
	if e != nil && e.Config.Method != "" {
		b := e.responseBody
		if b != nil {
			errs, err := compare.CompareJsonBody(e.Config.ResponseTmpls, *b, e.Config.ComparisonParams)
			if len(errs) > 0 {
				msg := ""
				for _, v := range errs {
					msg += v.Error() + "\n"
				}
				return &contract.TestError{
					Title:         "response body differs",
					Expected:      e.Config.ResponseTmpls,
					Actual:        *b,
					Message:       msg,
					OriginalError: fmt.Errorf("response body differs: %v", msg),
				}
			}
			if err != nil {
				return fmt.Errorf("compare json failed: %w", err)
			}
		}
	}
	if e != nil && e.Config.ResponseStatus != nil {
		if *e.Config.ResponseStatus != e.responseStatus {
			return &contract.TestError{
				Title: "response status differs",
				Message: fmt.Sprintf(
					"status expected %v, got %v",
					color.GreenString("%v", *e.Config.ResponseStatus),
					color.RedString("%v", e.responseStatus),
				),

				OriginalError: fmt.Errorf("response status differs"),
			}
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

func request(r RequestConfig, b *bytes.Buffer, host string) (*http.Request, error) {
	req, err := http.NewRequest(
		strings.ToUpper(r.Method),
		host+r.RequestURL+r.QueryParams,
		b,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range r.HeadersVal {
		if strings.ToLower(k) == "host" {
			req.Host = v
		} else {
			req.Header.Add(k, v)
		}
	}
	return req, nil
}

func actualRequestBody(req *http.Request) string {
	if req.Body != nil {
		reqBodyStream, _ := req.GetBody()
		reqBody, _ := ioutil.ReadAll(reqBodyStream)
		return string(reqBody)
	}
	return ""
}

func newCommonRequest(host string, r RequestConfig) (*http.Request, error) {
	body := []byte(r.RequestTmpl)
	req, err := request(r, bytes.NewBuffer(body), host)
	if err != nil {
		return nil, err
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
