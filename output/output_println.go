package output

import (
	"context"
	"fmt"
	"time"

	"github.com/ixpectus/declarate/contract"
)

var bar *Bar

type OutputPrintln struct {
	WithProgressBar bool
}

func (o *OutputPrintln) Log(message contract.Message) {
	if o.WithProgressBar {
		o.logWithProgressBar(message)
	} else {
		o.log(message)
	}
}

func (o *OutputPrintln) tickPoll(ctx context.Context, duration time.Duration) {
	tick := time.NewTicker(duration)
	defer tick.Stop()
	for true {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:

		}
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
