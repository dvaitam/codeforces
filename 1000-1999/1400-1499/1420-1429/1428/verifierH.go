package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testInput struct {
	n, m int
	vals []int
}

var verifierDir string

func init() {
	if _, file, _, ok := runtime.Caller(0); ok {
		verifierDir = filepath.Dir(file)
	} else {
		verifierDir = "."
	}
}

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func extractAnswer(output string, n int) ([]int, error) {
	fields := strings.Fields(output)
	var answer []int
	for i := 0; i < len(fields); i++ {
		if fields[i] != "2" {
			continue
		}
		if i+n >= len(fields)+1 {
			return nil, fmt.Errorf("command 2 missing numbers")
		}
		if i+n-1 >= len(fields) {
			return nil, fmt.Errorf("not enough integers after command 2")
		}
		cur := make([]int, n-1)
		ok := true
		for j := 0; j < n-1; j++ {
			val, err := strconv.Atoi(fields[i+1+j])
			if err != nil {
				ok = false
				break
			}
			cur[j] = val
		}
		if ok {
			answer = cur
		}
	}
	if answer == nil {
		return nil, fmt.Errorf("no final answer found in output:\n%s", output)
	}
	return answer, nil
}

func formatInput(tc testInput) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for _, v := range tc.vals {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	return sb.String()
}

func generateTests() []testInput {
	tests := []testInput{
		{n: 2, m: 2, vals: []int{1, 1, 1, 1}},
		{n: 3, m: 3, vals: []int{0, 1, 2, 3, 4, 5, 6, 7, 8}},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		n := rng.Intn(50) + 2
		m := rng.Intn(19) + 2
		k := n * m
		vals := make([]int, k)
		for i := 0; i < k; i++ {
			vals[i] = rng.Intn(k + 1)
		}
		if n-1 > 0 {
			idx := rng.Perm(k)
			maxVal := vals[idx[0]]
			for i := 1; i < k; i++ {
				if vals[idx[i]] > maxVal {
					maxVal = vals[idx[i]]
				}
			}
			bonus := maxVal + rng.Intn(k+1) + 1
			for i := 0; i < n-1 && i < k; i++ {
				vals[idx[i]] = bonus
			}
		}
		tests = append(tests, testInput{n: n, m: m, vals: vals})
	}
	tests = append(tests, testInput{
		n:    100,
		m:    20,
		vals: make([]int, 100*20),
	})
	for i := range tests[len(tests)-1].vals {
		tests[len(tests)-1].vals[i] = i % (tests[len(tests)-1].n * tests[len(tests)-1].m)
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]
	reference := filepath.Join(verifierDir, "1428H.go")
	tests := generateTests()
	for i, tc := range tests {
		input := formatInput(tc)
		expOut, err := runProgram(reference, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expOut) == strings.TrimSpace(gotOut) {
			continue
		}
		expAns, err := extractAnswer(expOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on case %d: %v\noutput:\n%s", i+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := extractAnswer(gotOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output parse error: %v\noutput:\n%s", i+1, err, gotOut)
			os.Exit(1)
		}
		if len(expAns) != len(gotAns) {
			fmt.Fprintf(os.Stderr, "case %d: expected %v got %v\n", i+1, expAns, gotAns)
			os.Exit(1)
		}
		match := true
		for j := range expAns {
			if expAns[j] != gotAns[j] {
				match = false
				break
			}
		}
		if !match {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", i+1, expAns, gotAns)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
