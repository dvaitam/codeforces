package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const refSource = "./1510D.go"

type testCase struct {
	n int
	d int
	a []int
}

type answer struct {
	neg bool
	seq []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refAns, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	refLogs := make([]float64, len(tests))
	for i, tc := range tests {
		if refAns[i].neg {
			refLogs[i] = math.Inf(-1)
			continue
		}
		val, err := evaluateSequence(tc, refAns[i].seq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference sequence invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		refLogs[i] = val
	}

	const tol = 1e-6
	for i, tc := range tests {
		exp := refAns[i]
		got := candAns[i]
		if exp.neg {
			if !got.neg {
				fmt.Fprintf(os.Stderr, "test %d: reference has no solution but candidate provided one\n", i+1)
				os.Exit(1)
			}
			continue
		}
		if got.neg {
			fmt.Fprintf(os.Stderr, "test %d: candidate output -1 but solution exists\n", i+1)
			os.Exit(1)
		}
		val, err := evaluateSequence(tc, got.seq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if math.Abs(val-refLogs[i]) > tol {
			fmt.Fprintf(os.Stderr, "test %d: suboptimal solution (log %.8f vs %.8f)\n", i+1, val, refLogs[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1510D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswers(output string, t int) ([]answer, error) {
	tokens := strings.Fields(output)
	ans := make([]answer, t)
	idx := 0
	for i := 0; i < t; i++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		if tokens[idx] == "-1" {
			ans[i] = answer{neg: true}
			idx++
			continue
		}
		k, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid length %q", i+1, tokens[idx])
		}
		idx++
		if k < 0 {
			return nil, fmt.Errorf("test %d: negative length", i+1)
		}
		if idx+k > len(tokens) {
			return nil, fmt.Errorf("test %d: insufficient sequence values", i+1)
		}
		seq := make([]int, k)
		for j := 0; j < k; j++ {
			val, err := strconv.Atoi(tokens[idx+j])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid sequence value %q", i+1, tokens[idx+j])
			}
			seq[j] = val
		}
		idx += k
		ans[i] = answer{seq: seq}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[idx])
	}
	return ans, nil
}

func evaluateSequence(tc testCase, seq []int) (float64, error) {
	if len(seq) == 0 {
		return 0, fmt.Errorf("sequence must be non-empty")
	}
	pos := 0
	for _, val := range seq {
		found := false
		for pos < tc.n {
			if tc.a[pos] == val {
				pos++
				found = true
				break
			}
			pos++
		}
		if !found {
			return 0, fmt.Errorf("sequence is not a subsequence of the input array")
		}
	}
	mod := 1
	sumLog := 0.0
	for _, val := range seq {
		mod = (mod * (val % 10)) % 10
		if val <= 0 {
			return 0, fmt.Errorf("invalid value %d in sequence", val)
		}
		sumLog += math.Log(float64(val))
	}
	if mod != tc.d {
		return 0, fmt.Errorf("product modulo 10 = %d, expected %d", mod, tc.d)
	}
	return sumLog, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(1510))
	var tests []testCase
	total := 0
	add := func(n, d int, arr []int) {
		if total+n > 500000 {
			return
		}
		cp := append([]int(nil), arr...)
		tests = append(tests, testCase{n: n, d: d, a: cp})
		total += n
	}

	add(4, 3, []int{3, 9, 9, 2})
	add(5, 3, []int{3, 2, 5, 6, 4})
	add(7, 1, []int{5, 9, 2, 6, 5, 4, 6})
	add(8, 2, []int{7, 1, 2, 6, 8, 3, 4, 5})
	add(4, 5, []int{3, 4, 5, 6})

	for total < 500000 {
		n := rng.Intn(2000) + 1
		if total+n > 500000 {
			n = 500000 - total
		}
		d := rng.Intn(10)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(1_000_000_000) + 1
		}
		add(n, d, arr)
		if len(tests) > 600 {
			break
		}
	}
	return tests
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.d)
		for i, v := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}
