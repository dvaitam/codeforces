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

const refSource = "2133F.go"

type testCase struct {
	n int
	e []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseRef(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseCand(candOut, refAns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		expectedK := refAns[i].k
		got := candAns[i]
		if expectedK == -1 {
			if got.k != -1 {
				fmt.Fprintf(os.Stderr, "test %d mismatch: expected -1 got %d detonations\n", i+1, got.k)
				fmt.Fprintf(os.Stderr, "n=%d e=%v\n", tc.n, tc.e)
				os.Exit(1)
			}
			continue
		}
		if got.k != expectedK {
			fmt.Fprintf(os.Stderr, "test %d mismatch in minimal detonations: expected %d got %d\n", i+1, expectedK, got.k)
			fmt.Fprintf(os.Stderr, "n=%d e=%v\n", tc.n, tc.e)
			os.Exit(1)
		}
		if err := verifySequence(tc, got.seq); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid sequence: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "n=%d e=%v\n", tc.n, tc.e)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	refPath, err := referencePath()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "ref_2133F_*.bin")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate verifier path")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, refSource), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

type refResult struct {
	k int
}

func parseRef(out string, tests int) ([]refResult, error) {
	lines := strings.Fields(out)
	idx := 0
	res := make([]refResult, tests)
	for t := 0; t < tests; t++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("not enough tokens for test %d", t+1)
		}
		k, err := strconv.Atoi(lines[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", lines[idx])
		}
		idx++
		res[t].k = k
		if k == -1 {
			continue
		}
		if k < 0 {
			return nil, fmt.Errorf("negative k")
		}
		if idx+k > len(lines) {
			return nil, fmt.Errorf("not enough tokens for sequence of test %d", t+1)
		}
		idx += k // skip sequence
	}
	if idx != len(lines) {
		return nil, fmt.Errorf("extra tokens in reference output")
	}
	return res, nil
}

type candResult struct {
	k   int
	seq []int
}

func parseCand(out string, ref []refResult) ([]candResult, error) {
	tokens := strings.Fields(out)
	idx := 0
	res := make([]candResult, len(ref))
	for t := 0; t < len(ref); t++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("not enough tokens for test %d", t+1)
		}
		k, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tokens[idx])
		}
		idx++
		if k == -1 {
			res[t] = candResult{k: -1}
			continue
		}
		if k <= 0 {
			return nil, fmt.Errorf("k must be positive for test %d", t+1)
		}
		if idx+k > len(tokens) {
			return nil, fmt.Errorf("not enough tokens for sequence of test %d", t+1)
		}
		seq := make([]int, k)
		for i := 0; i < k; i++ {
			val, err := strconv.Atoi(tokens[idx+i])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", tokens[idx+i])
			}
			seq[i] = val
		}
		idx += k
		res[t] = candResult{k: k, seq: seq}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens in candidate output")
	}
	return res, nil
}

func verifySequence(tc testCase, seq []int) error {
	n := tc.n
	e := tc.e
	if len(seq) == 0 {
		return fmt.Errorf("empty sequence")
	}
	next := make([]int, n+2)
	for i := 0; i <= n+1; i++ {
		next[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if next[x] != x {
			next[x] = find(next[x])
		}
		return next[x]
	}
	alive := n
	killRange := func(l, r int) {
		if l < 1 {
			l = 1
		}
		if r > n {
			r = n
		}
		x := find(l)
		for x <= r {
			alive--
			next[x] = x + 1
			x = find(x)
		}
	}

	for _, pos := range seq {
		if pos < 1 || pos > n {
			return fmt.Errorf("detonation index %d out of bounds", pos)
		}
		if find(pos) != pos {
			return fmt.Errorf("detonating already dead creeper at %d", pos)
		}
		if e[pos-1] == 0 {
			return fmt.Errorf("detonating creeper with zero power at %d", pos)
		}
		l := pos - e[pos-1] + 1
		r := pos + e[pos-1] - 1
		killRange(l, r)
		if alive == 0 {
			return nil
		}
	}
	if alive != 0 {
		return fmt.Errorf("creepers remain alive (%d)", alive)
	}
	return nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(e []int) {
		tests = append(tests, testCase{n: len(e), e: e})
	}

	// Fixed small cases including samples
	add([]int{0, 2, 2, 3, 0, 1})
	add([]int{0, 1, 2, 3})
	add([]int{1, 1, 1, 1})
	add([]int{0, 0, 0})
	add([]int{2, 4, 1, 3})
	add([]int{2, 0, 2, 4, 2, 2, 4, 1, 1})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	for len(tests) < 120 && totalN < 500000 {
		n := rng.Intn(3000) + 2 // 2..3001
		if totalN+n > 500000 {
			n = 500000 - totalN
		}
		e := make([]int, n)
		for i := 0; i < n; i++ {
			e[i] = rng.Intn(n + 1)
		}
		add(e)
		totalN += n
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.e {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
