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
		// msg := fmt.Sprintf("failed: %v:%v", message.Filename, message.Name)
		// messages = append(messages, colorFail.Sprint(msg))
		// pp.Println(messages)
		messageFormatted := strings.ReplaceAll(message.Message, "failed ", "")
		if message.Title != "" {
			messages = append(messages, message.Title+": \n"+messageFormatted)
		} else {
			messages = append(messages, prefix+messageFormatted)
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
	if message.Type == contract.MessageTypeNotify && !strings.Contains(message.Message, "start") {
		logText := colorNotify.Sprint(message.Message)
		messages = append(messages, logText)
	}
	if message.Type == contract.MessageTypeSuccess {
		logText := colorSuccess.Sprint(message.Message)
		messages = append(messages, prefix+logText)
	}

	return messages
}
