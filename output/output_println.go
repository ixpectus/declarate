package output

import (
	"fmt"

	"github.com/ixpectus/declarate/contract"
)

type OutputPrintln struct{}

func (o *OutputPrintln) Log(message contract.Message) {
	prefix := ""
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