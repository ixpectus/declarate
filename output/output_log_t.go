package output

import (
	"fmt"
	"log"
	"strings"

	"github.com/dailymotion/allure-go"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/tools"
)

func NewOutputT() *OutputT {
	return &OutputT{
		WithProgressBar: true,
	}
}

type OutputT struct {
	WithProgressBar bool
	report          contract.Report
}

func (o *OutputT) SetReport(r contract.Report) {
	o.report = r
}

func (o *OutputT) Log(message contract.Message) {
	if o.WithProgressBar {
		o.logWithProgressBar(message)
	} else {
		o.log(message)
	}
	if o.report != nil && message.Type == contract.MessageTypeError {
		o.report.Fail(errWrap(fmt.Sprintf("failed: %v:%v", tools.FilenameShort(message.Filename), message.Name)))
		o.report.AddAttachment("error details", allure.TextPlain, []byte(errMsgs("", message)))
	}
}

func (o *OutputT) logWithProgressBar(message contract.Message) {
	prefix := ""
	for i := 0; i < message.Lvl; i++ {
		prefix += " "
	}
	if message.Poll != nil {
		if bar == nil {
			if message.Type == contract.MessageTypeNotify {
				log.Println(prefix + message.Message)
			}
			bar = NewBar(message.Poll.Finish)
			go bar.Start()
		}
	} else {
		if message.PollResult != nil {
			if bar != nil {
				bar.SetCurrent()
			}
		}
		if bar != nil {
			bar.Stop()
			bar = nil
		}
		if message.Expected != "" {
			logText := colorNotify2.Sprint("expected: \n") + colorNotify.Sprint(message.Expected)
			log.Println(prefix + logText)
		}
		if message.Actual != "" {
			logText := colorNotify2.Sprint("got     : \n") + colorNotify.Sprint(message.Actual)
			log.Println(logText)
		}
		if message.Type == contract.MessageTypeError {
			msg := fmt.Sprintf("failed: %v", message.Name)
			if message.PollConditionFailed {
				msg = fmt.Sprintf("failed, poll condition: %v", message.Name)
			}
			log.Println(colorFail.Sprint(msg))
			if message.Title != "" {
				logText := colorFail.Sprint(message.Title) + ": \n" + message.Message
				log.Println(logText)
			} else {
				logText := colorFail.Sprint(message.Message)
				log.Println(prefix + logText)
			}

		}
		hideStart := strings.Contains(message.Message, "start") && !message.HasNestedSteps
		if message.Type == contract.MessageTypeNotify && !hideStart && message.Name != "" {
			logText := colorNotify.Sprint(message.Name)
			if strings.Contains(message.Message, "start") {
				logText = colorNotify.Sprint("start: `" + message.Name + "`")
			}
			log.Println(prefix + logText)
		}
		if message.Type == contract.MessageTypeSuccess && message.Name != "" {
			logText := colorSuccess.Sprint("passed: ") + colorNotify.Sprint("`"+message.Name+"`")
			log.Println(prefix + logText)
		}
	}
}

func (o *OutputT) log(message contract.Message) {
	msgs := messages(message)
	for _, v := range msgs {
		log.Println(v)
	}
}
