package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSourceC = "2000-2999/2100-2199/2160-2169/2167/2167C.go"

type testCaseC struct {
	n   int
	arr []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReferenceC()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTestsC()
	input := buildInputC(tests)

	refOut, err := runProgramC(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidateC(os.Args[1], input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refAns, err := parseOutputC(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutputC(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		refSeq := refAns[i]
		candSeq := candAns[i]
		if len(candSeq) != tc.n {
			fmt.Fprintf(os.Stderr, "test %d: expected %d numbers, got %d\n", i+1, tc.n, len(candSeq))
			os.Exit(1)
		}
		if !isLexicographicallyEqual(candSeq, refSeq) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d:\nexpected: %v\ngot:      %v\n", i+1, refSeq, candSeq)
			os.Exit(1)
		}
		if !isReachable(tc.arr, candSeq) {
			fmt.Fprintf(os.Stderr, "test %d: candidate sequence is not reachable under parity swaps\n", i+1)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReferenceC() (string, error) {
	tmp, err := os.CreateTemp("", "2167C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceC))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidateC(path, input string) (string, error) {
	cmd := commandForC(path)
	return runWithInputC(cmd, input)
}

func runProgramC(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInputC(cmd, input)
}

func commandForC(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInputC(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputC(output string, tests []testCaseC) ([][]int, error) {
	tokens := strings.Fields(output)
	ans := make([][]int, len(tests))
	pos := 0
	for i, tc := range tests {
		if pos+tc.n > len(tokens) {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		cur := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return nil, fmt.Errorf("test %d: token %q is not an integer", i+1, tokens[pos])
			}
			cur[j] = val
			pos++
		}
		ans[i] = cur
	}
	if pos != len(tokens) {
		if pos < len(tokens) {
			return nil, fmt.Errorf("extra output detected starting at token %q", tokens[pos])
		}
		return nil, fmt.Errorf("extra output detected")
	}
	return ans, nil
}

func isLexicographicallyEqual(a, b []int) bool {
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

func isReachable(original, candidate []int) bool {
	if len(original) != len(candidate) {
		return false
	}
	counts := make(map[int]int, len(original))
	for _, v := range original {
		counts[v]++
	}
	for _, v := range candidate {
		counts[v]--
		if counts[v] < 0 {
			return false
		}
	}
	for _, v := range counts {
		if v != 0 {
			return false
		}
	}
	return true
}

func generateTestsC() []testCaseC {
	const limit = 200000
	rng := rand.New(rand.NewSource(21672067))
	var tests []testCaseC
	total := 0

	add := func(arr []int) {
		if len(arr) == 0 {
			return
		}
		if total+len(arr) > limit {
			return
		}
		cpy := append([]int(nil), arr...)
		tests = append(tests, testCaseC{n: len(arr), arr: cpy})
		total += len(arr)
	}

	add([]int{4, 2, 3, 1})
	add([]int{3, 2, 1, 3, 4})
	add([]int{3, 7, 5, 1})
	add([]int{1000000000, 2})
	add([]int{1, 3, 5})
	add([]int{1, 2, 3, 5, 7})
	add([]int{2, 4, 8, 6})

	for i := 0; i < 50; i++ {
		n := 1 + rng.Intn(10)
		arr := randomArray(rng, n)
		add(arr)
	}
	for i := 0; i < 40; i++ {
		n := 10 + rng.Intn(100)
		arr := randomArray(rng, n)
		add(arr)
	}
	for i := 0; i < 10; i++ {
		n := 100 + rng.Intn(500)
		arr := randomArray(rng, n)
		add(arr)
	}
	if total < limit {
		arr := randomArray(rng, limit-total)
		add(arr)
	}

	return tests
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1_000_000_000) + 1
	}
	return arr
}

func buildInputC(tests []testCaseC) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		writeArrayC(&b, tc.arr)
	}
	return b.String()
}

func writeArrayC(b *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteByte('\n')
}
