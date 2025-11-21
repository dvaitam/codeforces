package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testInput struct {
	n   int
	arr []int
}

type testCase struct {
	desc  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, refOut)
			os.Exit(1)
		}
		exp, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		got, err := parseAnswers(out, tc.t)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		for idx := 0; idx < tc.t; idx++ {
			if exp[idx] != got[idx] {
				fmt.Printf("Wrong answer on test %d (%s) case %d: expected %d got %d\nInput:\n%s", i+1, tc.desc, idx+1, exp[idx], got[idx], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref2153C.bin"
	cmd := exec.Command("go", "build", "-o", path, "2153C.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseAnswers(out string, t int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to read answer %d: %v", i+1, err)
		}
	}
	if extra := strings.TrimSpace(readRemaining(reader)); extra != "" {
		return nil, fmt.Errorf("unexpected extra output: %q", extra)
	}
	return ans, nil
}

func readRemaining(r *bufio.Reader) string {
	var sb strings.Builder
	for {
		line, err := r.ReadString('\n')
		sb.WriteString(line)
		if err != nil {
			break
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, cases []testInput) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(cases))
		for _, c := range cases {
			fmt.Fprintf(&sb, "%d\n", c.n)
			for i, v := range c.arr {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{
			desc:  desc,
			input: sb.String(),
			t:     len(cases),
		})
	}

	add("example-like", []testInput{
		{n: 3, arr: []int{5, 5, 7}},
		{n: 3, arr: []int{4, 5, 7}},
		{n: 3, arr: []int{5, 5, 10}},
	})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 60 {
		numCases := rng.Intn(4) + 1
		cases := make([]testInput, numCases)
		totalN := 0
		for i := 0; i < numCases; i++ {
			n := rng.Intn(50) + 3
			totalN += n
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rng.Intn(30) + 1
			}
			cases[i] = testInput{n: n, arr: arr}
		}
		add(fmt.Sprintf("random-small-%d", len(tests)), cases)
	}

	// Stress tests
	large1 := testInput{
		n:   200000,
		arr: make([]int, 200000),
	}
	for i := range large1.arr {
		large1.arr[i] = (i % 1000) + 1
	}
	large2 := testInput{
		n:   200000,
		arr: make([]int, 200000),
	}
	for i := range large2.arr {
		if i%2 == 0 {
			large2.arr[i] = 1
		} else {
			large2.arr[i] = 1000000000
		}
	}
	add("large-cases", []testInput{large1, large2})

	return tests
}
