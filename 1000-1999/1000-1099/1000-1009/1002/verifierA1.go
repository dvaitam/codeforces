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

type Qubit struct{ ptr *int }

var valid map[*int]bool
var counts map[*int]int

func H(q Qubit) {
	if q.ptr == nil || !valid[q.ptr] {
		fmt.Println("invalid qubit")
		os.Exit(0)
	}
	counts[q.ptr]++
}

func main() {
	const n = %d
	tokens := make([]int, n)
	qs := make([]Qubit, n)
	valid = make(map[*int]bool, n)
	counts = make(map[*int]int, n)
	for i := 0; i < n; i++ {
		qs[i] = Qubit{ptr: &tokens[i]}
		valid[qs[i].ptr] = true
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic: %%v\n", r)
			os.Exit(0)
		}
	}()
	GenerateSuperposition(qs)
	for i := 0; i < n; i++ {
		if counts[qs[i].ptr] != 1 {
			fmt.Printf("qubit %%d called %%d times\n", i, counts[qs[i].ptr])
			return
		}
	}
	fmt.Println("OK")
}
`

func runCase(candidate string, n int) error {
	src := fmt.Sprintf(harnessTemplate, n)
	tmp, err := os.CreateTemp("", "harness-*.go")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(src); err != nil {
		tmp.Close()
		return fmt.Errorf("failed to write harness: %v", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("failed to close harness: %v", err)
	}

	cmd := exec.Command("go", "run", tmp.Name(), candidate)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/candidate.go")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	const maxN = 8
	for n := 1; n <= maxN; n++ {
		if err := runCase(candidate, n); err != nil {
			fmt.Fprintf(os.Stderr, "case n=%d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", maxN)
}
