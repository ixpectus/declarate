package run

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/ixpectus/declarate/contract"
)

func (r *Runner) logPoll(
	fileName string,
	v runConfig,
	pollInfo contract.PollInfo,
	d time.Duration,
	estimated time.Duration,
) {
	r.output.Log(contract.Message{
		Filename: r.filenameShort(fileName),
		Name:     v.Name,
		Poll:     &pollInfo,
		Message: fmt.Sprintf(
			"poll %s:%s, wait %v, estimated %v",
			r.filenameShort(fileName),
			v.Name,
			d,
			estimated.Truncate(time.Second),
		),
		Type: contract.MessageTypeNotify,
	})
}

func (r *Runner) logStart(fileName string, v runConfig, lvl int) {
	r.output.Log(contract.Message{
		Filename:       fileName,
		Name:           v.Name,
		HasNestedSteps: len(v.Steps) > 0,
		HasPoll:        len(v.Poll.PollInterval()) > 0,
		Message:        fmt.Sprintf("start %v:%v", fileName, v.Name),
		Type:           contract.MessageTypeNotify,
	})
}

func (r *Runner) logSkip(name, fileName string, lvl int) {
	r.output.Log(contract.Message{
		Filename:   fileName,
		Lvl:        lvl,
		Name:       name,
		ActionType: "skip",
		Message:    fmt.Sprintf("skipped for file %s: %s", r.filenameShort(fileName), name),
		Type:       contract.MessageTypeNotify,
	})
}

func (r *Runner) logPass(name, fileName string, res *Result, lvl int) {
	r.output.Log(contract.Message{
		Filename:   fileName,
		Lvl:        lvl,
		Name:       name,
		Message:    fmt.Sprintf("passed %v:%v", r.filenameShort(fileName), name),
		Type:       contract.MessageTypeSuccess,
		PollResult: res.PollResult,
	})
}

func (r *Runner) logRunFail(name, fileName string, err error, res *Result) {
	r.output.Log(contract.Message{
		Filename:            fileName,
		Name:                name,
		Message:             fmt.Sprintf("run failed for file %s: %s", r.filenameShort(fileName), err),
		Type:                contract.MessageTypeError,
		PollResult:          res.PollResult,
		PollConditionFailed: res.PollConditionFailed,
	})
}

func (r *Runner) logErr(res Result) {
	var errTest *contract.TestError
	if errors.As(res.Err, &errTest) {
		r.output.Log(contract.Message{
			Filename: res.FileName,
			Name:     res.Name,
			Message:  res.Err.Error(),
			Title: fmt.Sprintf(
				"failed %v:%v\n%v",
				res.FileName,
				res.Name,
				errTest.Title,
			),
			Expected:            errTest.Expected,
			Actual:              errTest.Actual,
			Lvl:                 res.Lvl,
			Type:                contract.MessageTypeError,
			PollResult:          res.PollResult,
			PollConditionFailed: res.PollConditionFailed,
		})
		return
	}
	if res.Err != nil {
		r.output.Log(contract.Message{
			Filename:            res.FileName,
			Name:                res.Name,
			Message:             fmt.Sprintf("failed %v", res.Err),
			Lvl:                 res.Lvl,
			Type:                contract.MessageTypeError,
			PollConditionFailed: res.PollConditionFailed,
		})
	}
}

func (r *Runner) filenameShort(fileName string) string {
	parts := strings.Split(fileName, "/")
	if len(parts) > 4 {
		return path.Base(fileName)
	}
	return fileName
}
