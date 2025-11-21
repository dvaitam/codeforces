package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2130-2139/2134/2134E.go"

type testCase struct {
	n int
	a []int
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseCandidate(out string, tests []testCase) ([][]int, error) {
	lines := strings.Split(out, "\n")
	results := make([][]int, len(tests))
	idx := 0
	for _, line := range lines {
		if idx >= len(tests) {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == "!" {
			fields = fields[1:]
		}
		if len(fields) != tests[idx].n {
			continue
		}
		arr := make([]int, tests[idx].n)
		ok := true
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				ok = false
				break
			}
			if v != 1 && v != 2 {
				ok = false
				break
			}
			arr[i] = v
		}
		if !ok {
			continue
		}
		results[idx] = arr
		idx++
	}
	if idx != len(tests) {
		return nil, fmt.Errorf("could not find answers for all test cases (found %d of %d)", idx, len(tests))
	}
	return results, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, a: []int{1, 2}},
		{n: 2, a: []int{2, 1}},
		{n: 3, a: []int{1, 1, 2}},
		{n: 3, a: []int{2, 2, 1}},
		{n: 4, a: []int{1, 2, 1, 2}},
	}
}

type fastRNG struct{ s uint64 }

func newRNG() *fastRNG {
	return &fastRNG{s: uint64(time.Now().UnixNano())}
}

func (r *fastRNG) Intn(n int) int {
	// xorshift64*
	r.s ^= r.s >> 12
	r.s ^= r.s << 25
	r.s ^= r.s >> 27
	return int((r.s * 2685821657736338717) % uint64(n))
}

func randomTests() []testCase {
	rng := newRNG()
	tests := []testCase{}
	total := 0
	for len(tests) < 30 && total < 900 {
		n := rng.Intn(20) + 2
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				a[i] = 1
			} else {
				a[i] = 2
			}
		}
		tests = append(tests, testCase{n: n, a: a})
		total += n
	}
	for len(tests) < 45 && total < 1000 {
		n := rng.Intn(100) + 50
		a := make([]int, n)
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				a[i] = 1
			} else {
				a[i] = 2
			}
		}
		tests = append(tests, testCase{n: n, a: a})
		total += n
	}
	return tests
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := deterministicTests()
	tests = append(tests, randomTests()...)

	if totalN(tests) > 1000 {
		tests = tests[:len(tests)-1]
	}

	input := buildInput(tests)

	// Build reference binary to ensure it compiles (not used for checking answers).
	refBin, err := func() (string, error) {
		tmp, err := os.CreateTemp("", "2134E-ref-*")
		if err != nil {
			return "", err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), nil
	}()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Remove(refBin)

	out, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ans, err := parseCandidate(out, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if len(ans[i]) != tc.n {
			fmt.Fprintf(os.Stderr, "test %d: expected %d values, got %d\n", i+1, tc.n, len(ans[i]))
			os.Exit(1)
		}
		for j := 0; j < tc.n; j++ {
			if ans[i][j] != tc.a[j] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at position %d: expected %d got %d\n", i+1, j+1, tc.a[j], ans[i][j])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d test cases passed\n", len(tests))
}
