package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "2000-2999/2000-2099/2040-2049/2044/2044E.go"

type testCase struct {
	name string
	data []query
}

type query struct {
	k, l1, r1, l2, r2 int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		input := buildInput(tc)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}
		refAns, err := parseOutput(len(tc.data), refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, refOut)
			os.Exit(1)
		}

		expected := evaluate(tc)
		if !equalAnswers(refAns, expected) {
			fmt.Fprintf(os.Stderr, "reference mismatch simulation on test %d (%s)\ninput:\n%soutput:\n%s",
				idx+1, tc.name, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(len(tc.data), candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, input, candOut)
			os.Exit(1)
		}
		if !equalAnswers(candAns, expected) {
			fmt.Fprintf(os.Stderr, "candidate mismatch on test %d (%s)\ninput:\n%soutput:\n%s",
				idx+1, tc.name, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2044E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2044E.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
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
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(expected int, output string) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(fields))
	}
	ans := make([]int64, expected)
	for i, token := range fields {
		val, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		ans[i] = val
	}
	return ans, nil
}

func equalAnswers(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.data)))
	for _, q := range tc.data {
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", q.k, q.l1, q.r1, q.l2, q.r2))
	}
	return sb.String()
}

func evaluate(tc testCase) []int64 {
	res := make([]int64, len(tc.data))
	for idx, q := range tc.data {
		res[idx] = solve(q.k, q.l1, q.r1, q.l2, q.r2)
	}
	return res
}

func solve(k, l1, r1, l2, r2 int64) int64 {
	ans := int64(0)
	p := int64(1)
	for {
		left := l1
		if x := ceilDiv(l2, p); x > left {
			left = x
		}
		right := r1
		if y := r2 / p; y < right {
			right = y
		}
		if left <= right {
			ans += right - left + 1
		}
		if p > r2/k {
			break
		}
		p *= k
	}
	return ans
}

func ceilDiv(a, b int64) int64 {
	return (a + b - 1) / b
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "sample",
			data: []query{
				{2, 2, 6, 2, 12},
				{2, 1, 1000000000, 1, 1000000000},
				{3, 5, 7, 15, 63},
				{1000000000, 1, 5, 6, 1000000000},
				{15, 17, 78, 2596, 20914861},
			},
		},
		{
			name: "edge_small",
			data: []query{
				{2, 1, 1, 1, 1},
				{2, 1, 2, 1, 4},
				{3, 1, 3, 1, 27},
			},
		},
	}

	rng := rand.New(rand.NewSource(123456789))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(10) + 1
	data := make([]query, t)
	for i := 0; i < t; i++ {
		k := rng.Int63n(1_000_000) + 2
		l1 := rng.Int63n(1_000_000_000) + 1
		r1 := l1 + rng.Int63n(1_000_000_000-l1+1)
		l2 := rng.Int63n(1_000_000_000) + 1
		r2 := l2 + rng.Int63n(1_000_000_000-l2+1)
		data[i] = query{k, l1, r1, l2, r2}
	}
	return testCase{
		name: fmt.Sprintf("random_%d", idx+1),
		data: data,
	}
}
