package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const numTestsA = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "binA*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, nil, nil
}

func run(bin string) (string, error) {
	cmd := exec.Command(bin)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func runCase(bin string) error {
	out, err := run(bin)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("expected empty output, got %q", out)
	}
	return nil
}

func main() {
	argIdx := 1
	if len(os.Args) >= 3 && os.Args[1] == "--" {
		argIdx = 2
	}
	if len(os.Args) != argIdx+1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin, cleanup, err := prepareBinary(os.Args[argIdx])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	for i := 0; i < numTestsA; i++ {
		if err := runCase(bin); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
