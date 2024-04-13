package script

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dailymotion/allure-go"
)

func (e *ScriptCmd) run(scriptPath string) (string, error) {
	cmds := strings.Split(strings.TrimRight(scriptPath, "\n"), " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)

	if e.Config.NoWait {
		if err := cmd.Start(); err != nil {
			return "", fmt.Errorf("cmd start: %w", err)
		}

		return "", nil
	}
	bb := bytes.Buffer{}
	errBB := bytes.Buffer{}
	cmd.Stdout = &bb
	cmd.Stderr = &errBB
	cmd.Env = os.Environ()
	if e.report != nil {
		e.report.AddAttachment("command", allure.TextPlain, []byte(cmd.String()))
	}
	if err := cmd.Run(); err != nil {
		log.Println(cmd.String())
		if e.report != nil {
			e.report.AddAttachment("stdout", allure.TextPlain, bb.Bytes())
			e.report.AddAttachment("stderr", allure.TextPlain, errBB.Bytes())
		}
		return "", fmt.Errorf("process finished with error = %v, output %v, std err %v", err, string(bb.Bytes()), string(errBB.Bytes()))
	}
	e.report.AddAttachment("stdout", allure.TextPlain, bb.Bytes())

	return bb.String(), nil
}
