package report

import (
	"testing"

	"github.com/dailymotion/allure-go"
)

type EmptyReport struct{}

func NewEmptyReport() *EmptyReport {
	return &EmptyReport{}
}

func (a *EmptyReport) Fail(err error) {
}

func (a *EmptyReport) AddAttachment(name string, mimeType allure.MimeType, content []byte) error {
	return nil
}

func (a *EmptyReport) Test(t *testing.T, action func(), options ReportOptions) {
	action()
}

func (a *EmptyReport) Action(action func()) allure.Option {
	return nil
}

func (a *EmptyReport) Description(description string) allure.Option {
	return nil
}

func (a *EmptyReport) Step(s ReportOptions, action func()) {
	action()
}
