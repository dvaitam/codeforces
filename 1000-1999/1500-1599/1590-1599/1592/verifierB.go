package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	n, x int
	a    []int
}

func parseInput(data []byte) ([]testCase, error) {
	reader := bytes.NewReader(data)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	tests := make([]testCase, t)
	total := 0
	for i := 0; i < t; i++ {
		var n, x int
		if _, err := fmt.Fscan(reader, &n, &x); err != nil {
			return nil, fmt.Errorf("failed to read n, x for test %d: %v", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &arr[j]); err != nil {
				return nil, fmt.Errorf("failed to read a[%d] of test %d: %v", j+1, i+1, err)
			}
		}
		tests[i] = testCase{n: n, x: x, a: arr}
		total += n
	}
	_ = total
	return tests, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func expectedResults(tests []testCase) []string {
	ans := make([]string, len(tests))
	for idx, tc := range tests {
		n, x := tc.n, tc.x
		a := tc.a
		if 2*x <= n {
			ans[idx] = "YES"
			continue
		}
		b := append([]int(nil), a...)
		sort.Ints(b)
		left := n - x
		right := x - 1
		ok := true
		for i := left; i <= right; i++ {
			if a[i] != b[i] {
				ok = false
				break
			}
		}
		if ok {
			ans[idx] = "YES"
		} else {
			ans[idx] = "NO"
		}
	}
	return ans
}

func checkOutput(output string, expected []string) error {
	fields := strings.Fields(output)
	if len(fields) < len(expected) {
		return fmt.Errorf("expected %d answers, got %d", len(expected), len(fields))
	}
	for i, exp := range expected {
		if !strings.EqualFold(fields[i], exp) {
			return fmt.Errorf("test %d: expected %s, got %s", i+1, exp, fields[i])
		}
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
		os.Exit(1)
	}
	tests, err := parseInput(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	expected := expectedResults(tests)
	out, err := runCandidate(args[0], string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := checkOutput(out, expected); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
