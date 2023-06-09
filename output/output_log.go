package output

import (
	"log"
	"strings"

	"github.com/fatih/color"
	"github.com/ixpectus/declarate/contract"
)

var (
	colorNotify  = color.New(color.FgCyan)
	colorNotify2 = color.New(color.FgHiCyan)
	colorSuccess = color.New(color.FgGreen)
	colorFail    = color.New(color.FgRed)
)

func New() *Output {
	return &Output{}
}

type Output struct{}

func (o *Output) Log(message contract.Message) {
	prefix := ""
	for i := 0; i < message.Lvl; i++ {
		prefix += " "
	}
	if message.Expected != "" {
		logText := colorNotify2.Sprint("expected: ") + colorNotify.Sprint(message.Expected)
		log.Println(prefix + logText)
	}
	if message.Actual != "" {
		logText := colorNotify2.Sprint("got     : ") + colorNotify.Sprint(message.Actual)
		log.Println(logText)
	}
	if message.Type == contract.MessageTypeError {
		if message.Title != "" {
			logText := colorFail.Sprint("FAILED: "+message.Title) + ": \n" + message.Message
			log.Println(logText)
		} else {
			logText := colorFail.Sprint("FAILED: " + message.Message)
			log.Println(prefix + logText)
		}
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
