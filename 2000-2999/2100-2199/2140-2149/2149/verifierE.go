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

type testCase struct {
	desc  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	path := "./ref2149E.bin"
	cmd := exec.Command("go", "build", "-o", path, "2149E.go")
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
			fmt.Fprintf(&sb, "%d %d %d %d\n", c.n, c.k, c.l, c.r)
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

	add("single-element", []testInput{{n: 1, k: 1, l: 1, r: 1, arr: []int{5}}})
	add("repeated", []testInput{{n: 4, k: 1, l: 1, r: 4, arr: []int{7, 7, 7, 7}}})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 60 {
		numCases := rng.Intn(4) + 1
		cases := make([]testInput, numCases)
		for i := 0; i < numCases; i++ {
			n := rng.Intn(40) + 1
			k := rng.Intn(n) + 1
			l := rng.Intn(n) + 1
			r := l + rng.Intn(n-l+1)
			arr := make([]int, n)
			for j := range arr {
				arr[j] = rng.Intn(5) + 1
			}
			cases[i] = testInput{n: n, k: k, l: l, r: r, arr: arr}
		}
		add(fmt.Sprintf("random-small-%d", len(tests)), cases)
	}

	large := []testInput{
		buildLarge(200000, 50000),
		buildLarge(200000, 100000),
	}
	add("large", large)

	return tests
}

type testInput struct {
	n   int
	k   int
	l   int
	r   int
	arr []int
}

func buildLarge(n, k int) testInput {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = (i % 100) + 1
	}
	return testInput{
		n:   n,
		k:   k,
		l:   1,
		r:   n,
		arr: arr,
	}
}
