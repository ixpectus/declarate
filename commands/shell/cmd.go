package shell

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dailymotion/allure-go"
)

func CmdGet(scriptPath string) *exec.Cmd {
	if strings.Contains(scriptPath, "\"") || strings.Contains(scriptPath, "|") {
		return CmdGetWithBash(scriptPath)
	}
	rawParts := strings.Split(scriptPath, " ")
	parts := make([]string, 0, len(rawParts))
	for _, v := range rawParts {
		if v != " " && v != "" {
			parts = append(parts, v)
		}
	}
	var cmd *exec.Cmd
	if len(parts) > 0 {
		cmd = exec.Command(parts[0], parts[1:]...)
	} else {
		cmd = exec.Command(parts[0])
	}
	cmd.Env = os.Environ()
	return cmd
}

func CmdGetWithBash(scriptPath string) *exec.Cmd {
	return exec.Command("bash", "-c", scriptPath)
}

func (e *ShellCmd) run(command string) ([]string, error) {
	commands := strings.Split(command, "\n")
	res := []string{}
	for _, v := range commands {
		v = strings.Trim(v, " ")
		if v == "" {
			continue
		}
		cmd := CmdGet(v)
		bb := bytes.Buffer{}
		errBB := bytes.Buffer{}
		cmd.Stdout = &bb
		cmd.Stderr = &errBB
		cmd.Env = os.Environ()
		if e.report != nil {
			e.report.AddAttachment("command", allure.TextPlain, []byte(cmd.String()))
		}
		if err := cmd.Run(); err != nil {
			if e.report != nil {
				e.report.AddAttachment("stdout", allure.TextPlain, bb.Bytes())
				e.report.AddAttachment("stderr", allure.TextPlain, errBB.Bytes())
			}
			log.Println(cmd.String())

			return nil, fmt.Errorf("process finished with error = %v, output %v, std err %v", err, bb.String(), errBB.String())
		}
		if e.report != nil {
			e.report.AddAttachment("stdout", allure.TextPlain, bb.Bytes())
		}
		res = append(res, bb.String())
	}

	return res, nil
}
