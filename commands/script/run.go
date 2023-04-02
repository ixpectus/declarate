package script

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Run(scriptPath string) (string, error) {
	cmds := strings.Split(strings.TrimRight(scriptPath, "\n"), " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	bb := bytes.Buffer{}
	errBB := bytes.Buffer{}
	cmd.Stdout = &bb
	cmd.Stderr = &errBB
	cmd.Env = os.Environ()
	log.Println(cmd.String())
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("process finished with error = %v, output %v, std err %v", err, string(bb.Bytes()), string(errBB.Bytes()))
	}

	return bb.String(), nil
}
