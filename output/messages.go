package output

import (
	"strings"

	"github.com/ixpectus/declarate/contract"
)

func messages(message contract.Message) []string {
	prefix := ""
	messages := []string{}
	for i := 0; i < message.Lvl; i++ {
		prefix += " "
	}
	if message.Type == contract.MessageTypeError {
		messageFormatted := strings.ReplaceAll(message.Message, "failed ", "")
		if message.Title != "" {
			messages = append(messages, colorFail.Sprint(message.Title)+": \n"+messageFormatted)
		} else {
			messages = append(messages, colorFail.Sprint(prefix)+messageFormatted)
		}
		messages = append(messages, "")
	}
	if message.Expected != "" {
		logText := colorNotify2.Sprint("expected response: \n") + colorNotify.Sprint(message.Expected)
		messages = append(messages, prefix+logText)
	}
	if message.Actual != "" {
		logText := colorNotify2.Sprint("actual response: \n") + colorNotify.Sprint(message.Actual)
		messages = append(messages, prefix+logText)
	}

	showMessage := !strings.Contains(message.Message, "start") || message.HasNestedSteps || message.HasPoll
	if message.Type == contract.MessageTypeNotify && showMessage {
		logText := colorNotify.Sprint(message.Message)
		messages = append(messages, prefix+logText)
	}
	if message.Type == contract.MessageTypeSuccess {
		logText := colorSuccess.Sprint(message.Message)
		messages = append(messages, prefix+logText)
	}

	return messages
}
