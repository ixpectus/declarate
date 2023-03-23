package shell

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CmdGet(scriptPath string) *exec.Cmd {
	if strings.Contains(scriptPath, "\"") || strings.Contains(scriptPath, "|") {
		return CmdGetWithBash(scriptPath)
	}
	parts := strings.Split(scriptPath, " ")
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

func Run(command string) ([]string, error) {
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
		log.Println(cmd.String())
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("process finished with error = %v, output %v, std err %v", err, string(bb.Bytes()), string(errBB.Bytes()))
		}

		res = append(res, bb.String())
	}

	return res, nil
}
