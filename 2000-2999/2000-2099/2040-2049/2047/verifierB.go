package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var fact [11]int64

func init() {
	fact[0] = 1
	for i := 1; i < len(fact); i++ {
		fact[i] = fact[i-1] * int64(i)
	}
}

type testCase struct {
	s string
}

func main() {
	args := os.Args[1:]
	if len(args) >= 1 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierB.go /path/to/solution")
		os.Exit(1)
	}
	target := args[0]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	testGroups := append(deterministicTests(), randomTests()...)

	for idx, tc := range testGroups {
		input := buildInput(tc)

		oracleOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, oracleOut)
			os.Exit(1)
		}
		candOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, input, candOut)
			os.Exit(1)
		}

		expected, err := parseOutput(oracleOut, len(tc))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s", idx+1, err, oracleOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut, len(tc))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s", idx+1, err, candOut)
			os.Exit(1)
		}

		for t := range expected {
			s := tc[t].s
			bestVal := optimalAfterOneOp(s)
			candVal := permCount([]byte(candAns[t]))
			if len(candAns[t]) != len(s) {
				fmt.Fprintf(os.Stderr, "test %d case %d: length mismatch (got %d, expected %d)\ninput: %s\n", idx+1, t+1, len(candAns[t]), len(s), s)
				os.Exit(1)
			}
			if !canReach(s, candAns[t]) {
				fmt.Fprintf(os.Stderr, "test %d case %d: candidate string %q not reachable from %q with one operation\n", idx+1, t+1, candAns[t], s)
				os.Exit(1)
			}
			if candVal != bestVal {
				fmt.Fprintf(os.Stderr, "test %d case %d: suboptimal permutation count %d (expected %d)\n", idx+1, t+1, candVal, bestVal)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d test groups passed.\n", len(testGroups))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier location")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2047B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2047B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func buildInput(tc []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tc)))
	sb.WriteByte('\n')
	for _, c := range tc {
		sb.WriteString(strconv.Itoa(len(c.s)))
		sb.WriteByte('\n')
		sb.WriteString(c.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != expected {
		return nil, fmt.Errorf("expected %d lines, got %d", expected, len(lines))
	}
	return lines, nil
}

func deterministicTests() [][]testCase {
	return [][]testCase{
		{{s: "a"}},
		{{s: "ab"}, {s: "abc"}},
		{{s: "xyz"}, {s: "aaaa"}, {s: "abcd"}},
	}
}

func randomTests() [][]testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	groups := make([][]testCase, 0, 40)
	for len(groups) < cap(groups) {
		tcCount := rng.Intn(5) + 1
		group := make([]testCase, tcCount)
		for i := 0; i < tcCount; i++ {
			n := rng.Intn(10) + 1
			b := make([]byte, n)
			for j := 0; j < n; j++ {
				b[j] = byte('a' + rng.Intn(26))
			}
			group[i] = testCase{s: string(b)}
		}
		groups = append(groups, group)
	}
	return groups
}

func optimalAfterOneOp(s string) int64 {
	b := []byte(s)
	n := len(b)
	best := int64(math.MaxInt64)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			b[i] = s[j]
			val := permCount(b)
			if val < best {
				best = val
			}
		}
		b[i] = s[i]
	}
	return best
}

func permCount(b []byte) int64 {
	freq := [26]int{}
	for _, ch := range b {
		freq[ch-'a']++
	}
	ans := fact[len(b)]
	for _, f := range freq {
		if f > 1 {
			ans /= fact[f]
		}
	}
	return ans
}

func canReach(orig, candidate string) bool {
	if len(orig) != len(candidate) {
		return false
	}
	diff := 0
	var idx int
	for i := 0; i < len(orig); i++ {
		if orig[i] != candidate[i] {
			diff++
			idx = i
		}
	}
	if diff == 0 {
		return true
	}
	if diff > 1 {
		return false
	}
	for j := 0; j < len(orig); j++ {
		if candidate[idx] == orig[j] {
			return true
		}
	}
	return false
}
