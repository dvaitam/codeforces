package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	input  string
	output string
}

var tests []testCase

func init() {
	for i := 0; i < 100; i++ {
		a1 := i * 2
		a2 := 200 - i*2
		a3 := i
		arr := []int{a1, a2, a3}
		sort.Ints(arr)
		med := arr[1]
		tests = append(tests, testCase{
			input:  fmt.Sprintf("%d %d %d\n", a1, a2, a3),
			output: fmt.Sprintf("%d", med),
		})
	}
}

func runTest(binary string, t testCase, idx int) error {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(t.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test %d runtime error: %v\n%s", idx, err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(t.output)
	if got != want {
		return fmt.Errorf("test %d failed: expected %q got %q", idx, want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, t := range tests {
		if err := runTest(bin, t, i+1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("Accepted")
}
