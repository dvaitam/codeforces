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
	refSource    = "2000-2999/2000-2099/2060-2069/2067/2067B.go"
	randomTests  = 120
	maxNSquared  = 900000 // keep within statement limit 1e6 with margin
	maxNPerCase  = 1000
	baseValueMax = 1000
)

type testCase struct {
	n   int
	arr []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}

	expect := parseOutputs(expectRaw, len(tests))
	got := parseOutputs(gotRaw, len(tests))

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at test %d: expected %s, got %s", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2067B-ref-*")
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

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+5)

	// Include some deterministic cases and examples.
	tests = append(tests,
		testCase{n: 2, arr: []int{1, 1}},
		testCase{n: 4, arr: []int{1, 1, 1, 1}},
		testCase{n: 4, arr: []int{1, 2, 3, 4}},
		testCase{n: 6, arr: []int{3, 3, 4, 5, 3, 3}},
	)

	used := 0
	for _, tc := range tests {
		used += tc.n * tc.n
	}
	for i := 0; i < randomTests && used < maxNSquared; i++ {
		remain := maxNSquared - used
		maxN := maxNPerCase
		if maxN*maxN > remain {
			// choose n so that n^2 <= remain
			maxN = intSqrt(remain)
		}
		if maxN < 2 {
			break
		}
		n := rng.Intn(maxN/2)*2 + 2 // even n in [2,maxN]
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(baseValueMax) + 1
			if arr[j] > n {
				// ensure some values exceed n; allowed by constraints
				arr[j] = rng.Intn(n) + 1
			}
		}
		tests = append(tests, testCase{n: n, arr: arr})
		used += n * n
	}

	return tests
}

func intSqrt(x int) int {
	r := int(1)
	for (r+1)*(r+1) <= x {
		r++
	}
	return r
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
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

func parseOutputs(out string, t int) []string {
	lines := strings.Fields(out)
	if len(lines) != t {
		fail("expected %d tokens, got %d", t, len(lines))
	}
	res := make([]string, t)
	for i, s := range lines {
		res[i] = strings.ToUpper(strings.TrimSpace(s))
	}
	return res
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
