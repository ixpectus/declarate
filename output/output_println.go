package output

import (
	"fmt"

	"github.com/dailymotion/allure-go"

	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/tools"
)

var bar *Bar

type OutputPrintln struct {
	WithProgressBar bool
	report          contract.Report
}

func (o *OutputPrintln) SetReport(r contract.Report) {
	o.report = r
}

func (o *OutputPrintln) Log(message contract.Message) {
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

func (o *OutputPrintln) logWithProgressBar(message contract.Message) {
	prefix := ""
	if message.Poll != nil {
		if bar == nil {
			if message.Type == contract.MessageTypeNotify {
				fmt.Println(prefix + message.Message)
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
		for i := 0; i < message.Lvl; i++ {
			prefix += " "
		}
		if message.Expected != "" {
			fmt.Println("expected: " + message.Expected)
		}
		if message.Actual != "" {
			fmt.Println("got     : " + message.Actual)
		}
		if message.Type == contract.MessageTypeError {
			if message.Title != "" {
				fmt.Println(message.Title + ": \n" + message.Message)
			} else {
				fmt.Println(prefix + message.Message)
			}
		}
		if message.Type == contract.MessageTypeNotify {
			fmt.Println(prefix + message.Message)
		}
		if message.Type == contract.MessageTypeSuccess {
			fmt.Println(prefix + message.Message)
		}
	}
}

func (o *OutputPrintln) log(message contract.Message) {
	msgs := messages(message)
	for _, v := range msgs {
		fmt.Println(v)
	}
}
