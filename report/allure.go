package report

import (
	"os"
	"testing"

	"github.com/dailymotion/allure-go"
)

type AllureReport struct {
	path string
}

func NewAllureReport(path string) *AllureReport {
	if path == "" {
		path = "./"
	}
	os.Setenv("ALLURE_RESULTS_PATH", path)
	return &AllureReport{
		path: path,
	}
}

func (a *AllureReport) Fail(err error) {
	allure.Fail(err)
}

func (a *AllureReport) AddAttachment(name string, mimeType allure.MimeType, content []byte) error {
	return allure.AddAttachment(name, mimeType, content)
}

func (a *AllureReport) Test(t *testing.T, action func(), options ReportOptions) {
	allure.Test(
		t,
		allure.Action(action),
		allure.Description(options.Description),
		allure.Epic(options.Epic),
		allure.ID(options.ID),
		allure.Suite(options.Suite),
		allure.SubSuite(options.SubSuite),
		allure.Tags(options.Tags...),
	)
}

func (a *AllureReport) Description(description string) allure.Option {
	return allure.Description(description)
}

func (a *AllureReport) Step(s ReportOptions, action func()) {
	allure.Step(allure.Action(action), allure.Description(s.Description))
}
