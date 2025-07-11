package echo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dailymotion/allure-go"
)

type mockVars struct{}

func (m mockVars) Apply(s string) string {
	return s + "_applied"
}

func (m mockVars) Set(k, val string) error {
	return nil
}

func (m mockVars) SetAll(m2 map[string]string) (map[string]string, error) {
	return m2, nil
}

func (m mockVars) Get(k string) string {
	return k + "_value"
}

func (m mockVars) SetPersistent(k, val string) error {
	return nil
}

type mockReport struct{}

func (m mockReport) AddAttachment(name string, mimeType allure.MimeType, content []byte) error {
	return nil
}

func TestEcho_Do(t *testing.T) {
	msg := "hello"
	cfg := &Config{Message: msg}
	e := &Echo{Config: cfg, Vars: mockVars{}}
	err := e.Do()
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if e.Config.Message != msg+"_applied" {
		t.Errorf("Do() did not apply vars: got %v", e.Config.Message)
	}
}

func TestEcho_ResponseBody(t *testing.T) {
	resp := "resp"
	e := &Echo{Config: &Config{Message: "msg", Response: &resp}}
	if got := e.ResponseBody(); got == nil || *got != resp {
		t.Errorf("ResponseBody() = %v, want %v", got, resp)
	}
	e2 := &Echo{Config: &Config{}}
	if got := e2.ResponseBody(); got != nil {
		t.Errorf("ResponseBody() = %v, want nil", got)
	}
}

func TestEcho_ExpectedResponse(t *testing.T) {
	msg := "expected"
	e := &Echo{Config: &Config{Message: msg}}
	if got := e.ExpectedResponse(); got != msg {
		t.Errorf("ExpectedResponse() = %v, want %v", got, msg)
	}
	e2 := &Echo{Config: &Config{}}
	if got := e2.ExpectedResponse(); got != "" {
		t.Errorf("ExpectedResponse() = %v, want empty string", got)
	}
}

func TestEcho_Check(t *testing.T) {
	msg := "foo"
	resp := "foo"
	e := &Echo{Config: &Config{Message: msg, Response: &resp}}
	if err := e.Check(); err != nil {
		t.Errorf("Check() error = %v, want nil", err)
	}
	badResp := "bar"
	e2 := &Echo{Config: &Config{Message: msg, Response: &badResp}}
	err := e2.Check()
	if err == nil || err.Error() != fmt.Sprintf("expected %s, got %s", badResp, msg) {
		t.Errorf("Check() error = %v, want error", err)
	}
}

func TestUnmarshaller_Build(t *testing.T) {
	u := &Unmarshaller{}
	// extendedConfig style
	extMsg := "ext"
	extResp := "resp"
	unmarshalExt := func(v interface{}) error {
		switch vv := v.(type) {
		case *extendedConfig:
			vv.Echo = &struct {
				Message  string  `yaml:"message,omitempty"`
				Response *string `yaml:"response,omitempty"`
			}{Message: extMsg, Response: &extResp}
			return nil
		case *config:
			return nil
		}
		return errors.New("unknown type")
	}
	doer, err := u.Build(unmarshalExt)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	echo, ok := doer.(*Echo)
	if !ok || echo.Config.Message != extMsg || echo.Config.Response == nil || *echo.Config.Response != extResp {
		t.Errorf("Build() extendedConfig failed: %+v", echo)
	}
	// short config style
	shortMsg := "short"
	shortResp := "resp2"
	unmarshalShort := func(v interface{}) error {
		switch vv := v.(type) {
		case *extendedConfig:
			return nil
		case *config:
			vv.Message = shortMsg
			vv.Response = &shortResp
			return nil
		}
		return errors.New("unknown type")
	}
	doer, err = u.Build(unmarshalShort)
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	echo, ok = doer.(*Echo)
	if !ok || echo.Config.Message != shortMsg || echo.Config.Response == nil || *echo.Config.Response != shortResp {
		t.Errorf("Build() config failed: %+v", echo)
	}
}

func TestConfig_isEmpty(t *testing.T) {
	c := &Config{}
	if !c.isEmpty() {
		t.Errorf("Config.isEmpty() = false, want true")
	}
	c.Message = "msg"
	if c.isEmpty() {
		t.Errorf("Config.isEmpty() = true, want false")
	}
}

func Test_config_isEmpty(t *testing.T) {
	c := &config{}
	if !c.isEmpty() {
		t.Errorf("config.isEmpty() = false, want true")
	}
	c.Message = "msg"
	if c.isEmpty() {
		t.Errorf("config.isEmpty() = true, want false")
	}
}

func Test_extendedConfig_isEmpty(t *testing.T) {
	var ex *extendedConfig
	if !ex.isEmpty() {
		t.Errorf("extendedConfig.isEmpty() = false, want true (nil)")
	}
	ex = &extendedConfig{}
	if !ex.isEmpty() {
		t.Errorf("extendedConfig.isEmpty() = false, want true (no Echo)")
	}
	ex.Echo = &struct {
		Message  string  `yaml:"message,omitempty"`
		Response *string `yaml:"response,omitempty"`
	}{}
	if !ex.isEmpty() {
		t.Errorf("extendedConfig.isEmpty() = false, want true (empty Message)")
	}
	ex.Echo.Message = "msg"
	if ex.isEmpty() {
		t.Errorf("extendedConfig.isEmpty() = true, want false")
	}
}

func TestEcho_Setters(t *testing.T) {
	e := &Echo{}
	vars := mockVars{}
	report := mockReport{}

	e.SetVars(vars)
	e.SetReport(report)

	if e.Vars == nil {
		t.Errorf("SetVars failed: Vars is nil")
	}
	if e.Report == nil {
		t.Errorf("SetReport failed: Report is nil")
	}
}
