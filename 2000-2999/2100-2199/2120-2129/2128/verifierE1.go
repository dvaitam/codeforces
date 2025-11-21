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
	"time"
)

const (
	refSource   = "2000-2999/2100-2199/2120-2129/2128/2128E1.go"
	maxTests    = 200
	totalNLimit = 4000
)

type testCase struct {
	n int
	k int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

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
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		expV := refAns[i].v
		ans := candAns[i]
		if ans.v != expV {
			fmt.Fprintf(os.Stderr, "test %d: wrong v, expected %d got %d\n", i+1, expV, ans.v)
			os.Exit(1)
		}
		if ans.v < 1 || ans.v > tc.n {
			fmt.Fprintf(os.Stderr, "test %d: v out of range [1,n]: %d\n", i+1, ans.v)
			os.Exit(1)
		}
		if ans.l < 1 || ans.r < ans.l || ans.r > tc.n {
			fmt.Fprintf(os.Stderr, "test %d: invalid segment [%d,%d]\n", i+1, ans.l, ans.r)
			os.Exit(1)
		}
		if ans.r-ans.l+1 < tc.k {
			fmt.Fprintf(os.Stderr, "test %d: segment too short (len %d, need >= %d)\n", i+1, ans.r-ans.l+1, tc.k)
			os.Exit(1)
		}
		if !isMedian(tc.a[ans.l-1:ans.r], ans.v) {
			fmt.Fprintf(os.Stderr, "test %d: v=%d is not a median of subarray\n", i+1, ans.v)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

type triple struct {
	v int
	l int
	r int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2128E1-ref-*")
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.k)
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

func parseOutput(out string, t int) ([]triple, error) {
	fields := strings.Fields(out)
	if len(fields) != 3*t {
		return nil, fmt.Errorf("expected %d integers, got %d", 3*t, len(fields))
	}
	res := make([]triple, t)
	idx := 0
	for i := 0; i < t; i++ {
		v, err := strconv.Atoi(fields[idx])
		if err != nil {
			return nil, fmt.Errorf("test %d: failed to parse v: %v", i+1, err)
		}
		l, err := strconv.Atoi(fields[idx+1])
		if err != nil {
			return nil, fmt.Errorf("test %d: failed to parse l: %v", i+1, err)
		}
		r, err := strconv.Atoi(fields[idx+2])
		if err != nil {
			return nil, fmt.Errorf("test %d: failed to parse r: %v", i+1, err)
		}
		idx += 3
		res[i] = triple{v: v, l: l, r: r}
	}
	return res, nil
}

func isMedian(arr []int, v int) bool {
	m := len(arr)
	need := (m + 1) / 2
	le := 0
	ge := 0
	for _, x := range arr {
		if x <= v {
			le++
		}
		if x >= v {
			ge++
		}
	}
	return le >= need && ge >= need
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumN := 0

	add := func(tc testCase) {
		if len(tests) >= maxTests || sumN+tc.n > totalNLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Small deterministic cases.
	add(testCase{n: 1, k: 1, a: []int{5}})
	add(testCase{n: 2, k: 1, a: []int{1, 2}})
	add(testCase{n: 3, k: 2, a: []int{4, 1, 2}})
	add(testCase{n: 4, k: 3, a: []int{4, 1, 2, 4}})

	for len(tests) < maxTests && sumN < totalNLimit {
		n := rng.Intn(50) + 1
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}
		k := rng.Intn(n) + 1
		a := make([]int, n)
		for i := 0; i < n; i++ {
			switch rng.Intn(5) {
			case 0:
				a[i] = 1
			case 1:
				a[i] = n
			default:
				a[i] = rng.Intn(n) + 1
			}
		}
		add(testCase{n: n, k: k, a: a})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 1, k: 1, a: []int{1}})
	}
	return tests
}
