package databases

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// executeShellCmd executes the provided command in the bash shell
func executeShellCmd(command []string) (string, error) {

	// Build the command
	cmd := exec.Command("/bin/bash", "-c", strings.Join(command, " "))

	// Store Stdout and Stderr
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error executing command: '%v', returned error '%v' and Stderr: '%s'", cmd.String(), err, strings.TrimSuffix(stderr.String(), "\n"))
	}

	return strings.TrimSuffix(out.String(), "\n"), nil

}
