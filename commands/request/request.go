package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"moul.io/http2curl"

	"github.com/dailymotion/allure-go"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/tools"
	"github.com/ixpectus/declarate/variables"
)

const (
	modeFull  = "full"
	modeLight = "light"
)

type Request struct {
	Config         *RequestConfig
	defaultConfig  DefaultConfig
	mode           string
	Vars           contract.Vars
	Host           string
	responseBody   *string
	responseStatus int
	comparer       contract.Comparer
	report         contract.ReportAttachement
}

type Unmarshaller struct {
	host          string
	comparer      contract.Comparer
	defaultConfig DefaultConfig
}

type Option func(*Unmarshaller)

func OptionDefaultRequestConfig(config DefaultConfig) Option {
	return func(c *Unmarshaller) {
		c.defaultConfig = config
	}
}

func NewUnmarshaller(
	host string,
	comparer contract.Comparer,
	opts ...Option,
) *Unmarshaller {
	u := &Unmarshaller{
		host:          host,
		comparer:      comparer,
		defaultConfig: DefaultConfig{},
	}
	for _, v := range opts {
		v(u)
	}
	return u
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfg := &RequestConfig{}
	if err := unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshall request: %w", err)
	}
	if cfg == nil {
		return nil, nil
	}
	if cfg.RequestURL == "" {
		return nil, nil
	}
	mode := modeLight
	if cfg.FullResponse != nil {
		mode = modeFull
	}
	return &Request{
		Config:        cfg,
		Host:          u.host,
		comparer:      u.comparer,
		mode:          mode,
		defaultConfig: u.defaultConfig,
	}, nil
}

type DefaultConfig struct {
	HeadersVal map[string]string `json:"headers" yaml:"headers"`
}
type RequestConfig struct {
	Method           string                 `json:"method" yaml:"method"`
	RequestTmpl      string                 `json:"request" yaml:"request"`
	Response         *string                `json:"response" yaml:"response"`
	ResponseStatus   *string                `json:"responseStatus" yaml:"responseStatus"`
	FullResponse     *string                `json:"fullResponse" yaml:"fullResponse"`
	HeadersVal       map[string]string      `json:"headers" yaml:"headers"`
	QueryParams      string                 `json:"query" yaml:"query"`
	CookiesVal       map[string]string      `json:"cookies" yaml:"cookies"`
	ComparisonParams contract.CompareParams `json:"comparisonParams" yaml:"comparisonParams"`
	RequestURL       string                 `json:"path" yaml:"path"`
}

func (e *Request) SetVars(vv contract.Vars) {
	e.Vars = vv
}

func (e *Request) SetReport(r contract.ReportAttachement) {
	e.report = r
}

func (e *Request) GetConfig() interface{} {
	return e.Config
}

func (e *Request) applyHeadersVal(headers map[string]string) map[string]string {
	res := map[string]string{}
	for k, v := range headers {
		k = e.Vars.Apply(k)
		v = e.Vars.Apply(v)
		res[k] = v
	}
	return res
}

func (e *Request) IsValid() error {
	if (e.Config.Response != nil || e.Config.ResponseStatus != nil) && e.Config.FullResponse != nil {
		return fmt.Errorf("impossible to fill response and fullResponse simultaneously, choose one of")
	}
	if e.Config.Response != nil {
		c := *e.Config.Response
		c = variables.VariableRx.ReplaceAllString(c, "2")

		valid := json.Valid([]byte(c))
		if !valid {
			return fmt.Errorf("cannot parse response: `%v`", *e.Config.Response)
		}
	}
	if e.Config.FullResponse != nil {
		valid := json.Valid([]byte(*e.Config.FullResponse))
		if !valid {
			return fmt.Errorf("cannot parse response: `%v`", *e.Config.FullResponse)
		}
	}

	return nil
}

func (e *Request) Do() error {
	if e.Config.Method == "" {
		return nil
	}
	e.Config.QueryParams = e.Vars.Apply(e.Config.QueryParams)
	e.Config.RequestTmpl = e.Vars.Apply(e.Config.RequestTmpl)

	if e.Config.Response != nil {
		s := e.Vars.Apply(*e.Config.Response)
		s = strings.TrimSuffix(s, "\n")
		e.Config.Response = &s
	}
	if e.Config.FullResponse != nil {
		s := e.Vars.Apply(*e.Config.FullResponse)
		s = strings.TrimSuffix(s, "\n")
		e.Config.FullResponse = &s
	}
	if e.Config.ResponseStatus != nil {
		s := e.Vars.Apply(*e.Config.ResponseStatus)
		s = strings.TrimSuffix(s, "\n")
		e.Config.ResponseStatus = &s
	}
	e.Config.RequestURL = e.Vars.Apply(e.Config.RequestURL)

	defaultHeaders := e.applyHeadersVal(e.defaultConfig.HeadersVal)
	if defaultHeaders == nil {
		defaultHeaders = map[string]string{}
	}
	headers := e.applyHeadersVal(e.Config.HeadersVal)
	for k, v := range headers {
		defaultHeaders[k] = v
	}
	config := *e.Config
	config.HeadersVal = defaultHeaders
	req, err := newCommonRequest(e.Host, config)
	if err != nil {
		return err
	}
	client := &http.Client{}
	curlReq, _ := http2curl.GetCurlCommand(req)
	reqStart := time.Now()
	e.report.AddAttachment(fmt.Sprintf("request"), allure.TextPlain, []byte(curlReq.String()))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	reqFinish := time.Now()
	reqDuration := reqFinish.Sub(reqStart)
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return err
	}
	s := string(body)
	ss, _ := json.Marshal(resp.Header)
	r := fmt.Sprintf(`{"body":%v, "status":%v, "header": %s}`, s, resp.StatusCode, ss)
	if s == "" {
		r = fmt.Sprintf(`{"body":"", "status":%v, "header": %s}`, resp.StatusCode, ss)
	}

	if e.report != nil {
		e.report.AddAttachment("response", allure.ApplicationJson, []byte(tools.JSONPrettyPrint(r)))
		e.report.AddAttachment("meta", allure.TextPlain, []byte(tools.FormatVariables(map[string]string{
			"start":    reqStart.Format(time.RFC822),
			"finish":   reqFinish.Format(time.RFC822),
			"duration": reqDuration.Round(time.Millisecond).String(),
		})))
	}
	e.responseBody = tools.To(r)

	return nil
}

func (e *Request) ResponseBody() *string {
	if e.mode == modeFull {
		return e.responseBody
	}
	if e.responseBody != nil {
		b := gjson.Get(*e.responseBody, "body").String()
		return &b
	}

	return e.responseBody
}

func (e *Request) Check() error {
	if e.mode == modeFull {
		return e.checkFull()
	}

	return e.checkLight()
}

func (e *Request) checkLight() error {
	if e != nil && e.Config.Method != "" {
		b := e.responseBody
		if b != nil && e.Config.ResponseStatus != nil {
			gotStatus := gjson.Get(*b, "status")
			errs := e.comparer.Compare(
				*e.Config.ResponseStatus,
				gotStatus.String(),
				e.Config.ComparisonParams,
			)
			if len(errs) > 0 {
				msg := ""
				for i, v := range errs {
					if i < len(errs)-1 {
						msg += v.Error() + "\n"
					} else {
						msg += v.Error()
					}
				}

				return &contract.TestError{
					Title:         "response status differs",
					Expected:      *e.Config.ResponseStatus,
					Actual:        gotStatus.String(),
					Message:       msg,
					OriginalError: fmt.Errorf("response status differs: %v", msg),
				}
			}
		}
		if b != nil && e.Config.Response != nil {
			body := gjson.Get(*b, "body")
			expectedBody := e.Config.Response
			errs, err := e.comparer.CompareJsonBody(*expectedBody, body.String(), e.Config.ComparisonParams)
			realResponse := body.String()

			if len(errs) > 0 {
				msg := ""
				for i, v := range errs {
					if i < len(errs)-1 {
						msg += v.Error() + "\n"
					} else {
						msg += v.Error()
					}
				}

				expectedRemarshal, err := tools.JSONRemarshal(*expectedBody)
				if err != nil {
					return fmt.Errorf("remarshall expected response: %w", err)
				}

				actualRemarshal, err := tools.JSONRemarshal(realResponse)
				if err != nil {
					return fmt.Errorf("remarshall real response %v: %w", realResponse, err)
				}
				return &contract.TestError{
					Title:         "response differs",
					Expected:      expectedRemarshal,
					Actual:        actualRemarshal,
					Message:       msg,
					OriginalError: fmt.Errorf("response body differs: %v", msg),
				}
			}
			if err != nil {
				return fmt.Errorf("compare json failed: %w", err)
			}
		}
	}
	return nil
}

func (e *Request) checkFull() error {
	if e != nil && e.Config.Method != "" {
		b := e.responseBody
		if b != nil && e.Config.FullResponse != nil {
			body := gjson.Get(*b, "body")
			gotStatus := gjson.Get(*b, "status")
			expectedBody := gjson.Get(*e.Config.FullResponse, "body")
			errs, err := e.comparer.CompareJsonBody(expectedBody.String(), body.String(), e.Config.ComparisonParams)
			realResponse := body.String()

			expectedResponse := expectedBody.String()
			if strings.Contains(*e.Config.FullResponse, "status") {
				status := gjson.Get(*e.Config.FullResponse, "status")
				compareErrs := e.comparer.Compare(status.Value(), gotStatus.Value(), contract.CompareParams{})

				if len(compareErrs) > 0 {
					realResponse = fmt.Sprintf(`{"body":%v, "status":%v}`, body, gotStatus.String())
					if status.Int() > 0 {
						expectedResponse = fmt.Sprintf(`{"body":%v, "status":%v}`, expectedResponse, status.String())
					} else {
						expectedResponse = fmt.Sprintf(`{"body":%v, "status":"%v"}`, expectedResponse, status.String())
					}
					errs = append(errs, fmt.Errorf("status differs, expected %s, got %s", status.String(), gotStatus.String()))
				}
			}

			if len(errs) > 0 {
				msg := ""
				for i, v := range errs {
					if i < len(errs)-1 {
						msg += v.Error() + "\n"
					} else {
						msg += v.Error()
					}
				}

				expectedRemarshal, err := tools.JSONRemarshal(expectedResponse)
				if err != nil {
					return fmt.Errorf("remarshall expected response: %w", err)
				}

				actualRemarshal, err := tools.JSONRemarshal(realResponse)
				if err != nil {
					return fmt.Errorf("remarshall real response %v: %w", realResponse, err)
				}
				return &contract.TestError{
					Title:         "response differs",
					Expected:      expectedRemarshal,
					Actual:        actualRemarshal,
					Message:       msg,
					OriginalError: fmt.Errorf("response body differs: %v", msg),
				}
			}
			if err != nil {
				return fmt.Errorf("compare json failed: %w", err)
			}
		}
	}
	return nil
}

func request(r RequestConfig, b *bytes.Buffer, host string) (*http.Request, error) {
	req, err := http.NewRequest(
		strings.ToUpper(r.Method),
		host+r.RequestURL+r.QueryParams,
		b,
	)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
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
