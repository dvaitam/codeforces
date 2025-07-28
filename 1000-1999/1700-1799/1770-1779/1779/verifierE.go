package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const expectOutput = "Problem E is interactive and cannot be automatically solved."

func run(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		out, err := run(bin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "run failed: %v", err)
			os.Exit(1)
		}
		if out != expectOutput {
			fmt.Fprintf(os.Stderr, "expected %q got %q", expectOutput, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
