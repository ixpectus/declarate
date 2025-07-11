package contract

import (
	"testing"

	"github.com/dailymotion/allure-go"

	"github.com/ixpectus/declarate/report"
)

type Report interface {
	Fail(err error)
	AddAttachment(name string, mimeType allure.MimeType, content []byte) error
	Test(t *testing.T, action func(), options report.ReportOptions)
	Step(s report.ReportOptions, action func())
}

type ReportAttachement interface {
	AddAttachment(name string, mimeType allure.MimeType, content []byte) error
}
