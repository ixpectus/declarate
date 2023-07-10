package run

import "github.com/ixpectus/declarate/contract"

type Result struct {
	Err        error
	Name       string
	Lvl        int
	FileName   string
	Response   *string
	PollResult *contract.PollResult
}
