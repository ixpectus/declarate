package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ixpectus/declarate/compare"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/tools"
	"github.com/tidwall/gjson"
)

type Request struct {
	Config         *RequestConfig
	defaultConfig  DefaultConfig
	Vars           contract.Vars
	Host           string
	responseBody   *string
	responseStatus int
	comparer       contract.Comparer
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
	u := &Unmarshaller{host: host, comparer: comparer, defaultConfig: DefaultConfig{}}
	for _, v := range opts {
		v(u)
	}
	return u
}

func (u *Unmarshaller) Build(unmarshal func(interface{}) error) (contract.Doer, error) {
	cfg := &RequestConfig{}
	if err := unmarshal(cfg); err != nil {
		return nil, err
	}

	if cfg == nil {
		return nil, nil
	}
	if cfg.RequestURL == "" {
		return nil, nil
	}
	return &Request{
		Config:        cfg,
		Host:          u.host,
		comparer:      u.comparer,
		defaultConfig: u.defaultConfig,
	}, nil
}

type DefaultConfig struct {
	HeadersVal map[string]string `json:"headers" yaml:"headers"`
}
type RequestConfig struct {
	Method           string                    `json:"method" yaml:"method"`
	RequestTmpl      string                    `json:"request" yaml:"request"`
	Response         map[int]string            `json:"response" yaml:"response"`
	ResponseFull     *string                   `json:"responseFull" yaml:"responseFull"`
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

func (e *Request) GetConfig() interface{} {
	return e.Config
}

func (e *Request) applyHeadersVal(headers map[string]string) map[string]string {
	for k, v := range headers {
		k = e.Vars.Apply(k)
		v = e.Vars.Apply(v)
		headers[k] = v
	}
	return headers
}

func (e *Request) IsValid() error {
	if e.Config.ResponseFull != nil && e.Config.Response != nil {
		return fmt.Errorf("only onf of response and responseFull should be set")
	}
	if e.Config.ResponseFull != nil {
		valid := json.Valid([]byte(*e.Config.ResponseFull))
		if !valid {
			return fmt.Errorf("cannot parse response_full: `%v`", *e.Config.ResponseFull)
		}
	}
	if e.Config.Response != nil {
		for _, v := range e.Config.Response {
			valid := json.Valid([]byte(v))
			if !valid {
				return fmt.Errorf("cannot parse response: `%v`", v)
			}
		}
	}
	return nil
}

func (e *Request) Do() error {
	if e.Config.Method != "" {
		e.Config.QueryParams = e.Vars.Apply(e.Config.QueryParams)
		e.Config.RequestTmpl = e.Vars.Apply(e.Config.RequestTmpl)
		if e.Config.Response != nil {
			for k, v := range e.Config.Response {
				r := fmt.Sprintf(`{"body":%v, "status":%v}`, v, k)
				e.Config.ResponseFull = &r
			}
		}
		if e.Config.ResponseFull != nil {
			s := e.Vars.Apply(*e.Config.ResponseFull)
			s = strings.TrimSuffix(s, "\n")
			e.Config.ResponseFull = &s
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
		// curlReq, _ := http2curl.GetCurlCommand(req)
		// pp.Println(curlReq)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return err
		}
		s := string(body)
		ss, _ := json.Marshal(resp.Header)
		r := fmt.Sprintf(`{"body":%v, "status":%v, "header": %s}`, s, resp.StatusCode, ss)
		e.responseBody = tools.To(r)
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
		if b != nil && e.Config.ResponseFull != nil {
			body := gjson.Get(*b, "body")
			gotStatus := gjson.Get(*b, "status")
			expectedBody := gjson.Get(*e.Config.ResponseFull, "body")
			errs, err := e.comparer.CompareJsonBody(expectedBody.String(), body.String(), e.Config.ComparisonParams)
			realResponse := body.String()
			expectedResponse := expectedBody.String()
			if strings.Contains(*e.Config.ResponseFull, "status") {
				status := gjson.Get(*e.Config.ResponseFull, "status")
				if status.String() != gotStatus.String() {
					realResponse = fmt.Sprintf(`{"body":%v, "status":%v}`, body, status.String())
					expectedResponse = fmt.Sprintf(`{"body":%v, "status":%v}`, expectedResponse, gotStatus.String())
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
					return err
				}
				actualRemarshal, err := tools.JSONRemarshal(realResponse)
				if err != nil {
					return err
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

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}
