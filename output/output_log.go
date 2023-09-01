package output

import (
	"fmt"
	"log"
	"strings"

	"github.com/dailymotion/allure-go"
	"github.com/fatih/color"
	"github.com/ixpectus/declarate/contract"
	"github.com/ixpectus/declarate/tools"
)

var (
	colorNotify  = color.New(color.FgCyan)
	colorNotify2 = color.New(color.FgHiCyan)
	colorSuccess = color.New(color.FgGreen)
	colorFail    = color.New(color.FgRed)
)

func New() *Output {
	return &Output{
		WithProgressBar: true,
	}
}

type Output struct {
	WithProgressBar bool
}

func (o *Output) Log(message contract.Message) {
	if o.WithProgressBar {
		o.logWithProgressBar(message)
	} else {
		o.log(message)
	}
}

func (o *Output) logWithProgressBar(message contract.Message) {
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
			msg := fmt.Sprintf("failed: %v:%v", tools.FilenameShort(message.Filename), message.Name)
			if message.PollConditionFailed {
				msg = fmt.Sprintf("failed, poll condition: %v:%v", tools.FilenameShort(message.Filename), message.Name)
			}
			log.Println(colorFail.Sprint(msg))
			if message.Title != "" {
				logText := colorFail.Sprint(message.Title) + ": \n" + message.Message
				log.Println(logText)
			} else {
				logText := colorFail.Sprint(message.Message)
				log.Println(prefix + logText)
			}
			allure.Fail(errWrap(msg))
			allure.AddAttachment("error details", allure.TextPlain, []byte(o.errMsgs(prefix, message)))

		}
		if message.Type == contract.MessageTypeNotify && !strings.Contains(message.Message, "start") {
			logText := colorNotify.Sprint(message.Message)
			log.Println(prefix + logText)
		}
		if message.Type == contract.MessageTypeSuccess {
			logText := colorSuccess.Sprint(message.Message)
			log.Println(prefix + logText)
		}
	}
}

func (o *Output) errMsgs(prefix string, message contract.Message) string {
	res := []string{}
	if message.Type == contract.MessageTypeError {
		logText := message.Message
		res = append(res, prefix+logText)
	}
	if message.Expected != "" {
		res = append(res, "expected: \n"+message.Expected)
	}
	if message.Actual != "" {
		res = append(res, "got     : \n"+message.Actual)
	}
	return strings.Join(res, "\n")
}

type errString struct {
	v string
}

func (e errString) Error() string {
	return e.v
}

func errWrap(s string) error {
	return errString{v: s}
}

func (o *Output) log(message contract.Message) {
	msgs := messages(message)
	for _, v := range msgs {
		log.Println(v)
	}
}
