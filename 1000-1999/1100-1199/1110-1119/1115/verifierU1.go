package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const harnessTemplate = `package main

import (
	"fmt"
	"os"
)

type harnessQubit struct {
	idx   int
	flips int
}

func (q *harnessQubit) X() {
	q.flips++
}

func main() {
	lengths := []int{%s}
	for _, n := range lengths {
		backing := make([]*harnessQubit, n)
		qs := make([]Qubit, n)
		for i := 0; i < n; i++ {
			backing[i] = &harnessQubit{idx: i}
			qs[i] = backing[i]
		}
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("panic on n=%d: %%v\n", n, r)
					os.Exit(0)
				}
			}()
			AntiDiagonal(qs)
		}()
		for i := 0; i < n; i++ {
			if backing[i].flips%%2 != 1 {
				fmt.Printf("n=%d qubit %d received %d X calls\n", n, i, backing[i].flips)
				return
			}
		}
	}
	fmt.Println("OK")
}
`

func runHarness(candidate string, lengths []int) error {
	var parts []string
	for _, v := range lengths {
		parts = append(parts, fmt.Sprintf("%d", v))
	}
	src := fmt.Sprintf(harnessTemplate, strings.Join(parts, ","))

	dir := filepath.Dir(candidate)
	tmpFile, err := os.CreateTemp(dir, "u1-harness-*.go")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(src); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write harness: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close harness: %v", err)
	}

	cmd := exec.Command("go", "run", tmpFile.Name(), candidate)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go run failed: %v\n%s", err, stderr.String())
	}
	output := strings.TrimSpace(stdout.String())
	if output != "OK" {
		return fmt.Errorf("harness reported failure:\n%s", stdout.String())
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierU1.go /path/to/candidate.go")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	lengths := []int{2, 3, 4, 5}
	if err := runHarness(candidate, lengths); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
