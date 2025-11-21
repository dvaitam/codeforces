package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA3.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s\n", err, refOut)
		os.Exit(1)
	}
	if err := ensureSilent(refOut); err != nil {
		fmt.Fprintf(os.Stderr, "reference violated silence requirement: %v\n", err)
		os.Exit(1)
	}

	out, err := runProgram(target)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n%s\n", err, out)
		os.Exit(1)
	}
	if err := ensureSilent(out); err != nil {
		fmt.Fprintf(os.Stderr, "wrong answer: %v\nOutput:\n%s\n", err, out)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}

func buildReference() (string, error) {
	path := "./ref1002A3.bin"
	cmd := exec.Command("go", "build", "-o", path, "1002A3.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(bin string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func ensureSilent(out string) error {
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("expected no output, got %q", truncate(out, 80))
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
