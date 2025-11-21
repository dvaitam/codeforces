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

const (
	refSource = "2000-2999/2000-2099/2060-2069/2069/2069D.go"
	limitSum  = 200000
)

type testCase struct {
	s string
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

	expectOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, expectOut)
		os.Exit(1)
	}

	gotOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, gotOut)
		os.Exit(1)
	}

	expect, err := parseOutputs(expectOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}

	got, err := parseOutputs(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expect[i] != got[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d (s=%s)\n", i+1, expect[i], got[i], tests[i].s)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2069D-ref-*")
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

func parseOutputs(output string, t int) ([]int, error) {
	tokens := strings.Fields(output)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[t])
	}
	ans := make([]int, t)
	for i := 0; i < t; i++ {
		val, err := strconv.Atoi(tokens[i])
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", tokens[i])
		}
		if val < 0 {
			return nil, fmt.Errorf("answer %d is negative", val)
		}
		ans[i] = val
	}
	return ans, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2069))
	var tests []testCase
	total := 0

	add := func(s string) {
		if total+len(s) > limitSum {
			return
		}
		tests = append(tests, testCase{s: s})
		total += len(s)
	}

	add("aa")
	add("abba")
	add("baab")
	add("abccbaabccba")
	add(makePalindrome("abcd"))
	add(makePairString(rng, 2))
	add(makePairString(rng, 4))
	add(makePairString(rng, 6))

	for total < limitSum && len(tests) < 200 {
		n := 2 + rng.Intn(200)
		if n%2 == 1 {
			n++
		}
		add(makePairString(rng, n))
		n2 := 2 + rng.Intn(200)
		if n2%2 == 1 {
			n2++
		}
		add(makeAlmostPalindrome(rng, n2))
	}

	for total < limitSum {
		n := min(2000, limitSum-total)
		if n%2 == 1 {
			n--
		}
		if n < 2 {
			break
		}
		s := makePairString(rng, n)
		add(s)
	}

	return tests
}

func makePairString(rng *rand.Rand, n int) string {
	if n%2 == 1 {
		n++
	}
	if n < 2 {
		n = 2
	}
	bytesArr := make([]byte, n)
	for i := 0; i < n/2; i++ {
		ch := byte('a' + rng.Intn(26))
		bytesArr[2*i] = ch
		bytesArr[2*i+1] = ch
	}
	shuffleBytes(rng, bytesArr)
	return string(bytesArr)
}

func makeAlmostPalindrome(rng *rand.Rand, n int) string {
	if n%2 == 1 {
		n++
	}
	half := n / 2
	left := make([]byte, half)
	for i := 0; i < half; i++ {
		left[i] = byte('a' + rng.Intn(26))
	}
	full := make([]byte, n)
	for i := 0; i < half; i++ {
		full[i] = left[i]
		full[n-1-i] = left[i]
	}
	// introduce some shuffle in the middle portion but keep counts even
	midLen := rng.Intn(n/2 + 1)
	if midLen%2 == 1 {
		midLen++
	}
	if midLen > n {
		midLen = n
	}
	start := rng.Intn(n - midLen + 1)
	sub := make([]byte, midLen)
	copy(sub, full[start:start+midLen])
	shuffleBytes(rng, sub)
	copy(full[start:start+midLen], sub)
	return string(full)
}

func makePalindrome(seed string) string {
	n := len(seed)
	res := make([]byte, 2*n)
	for i := 0; i < n; i++ {
		res[i] = seed[i]
		res[2*n-1-i] = seed[i]
	}
	return string(res)
}

func shuffleBytes(rng *rand.Rand, arr []byte) {
	rng.Shuffle(len(arr), func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&b, tc.s)
	}
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
