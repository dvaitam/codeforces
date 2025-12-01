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

// refSource points to the reference solution in this directory so Go does not
// try to resolve a module-like path under GOPATH.
const refSource = "193C.go"

type testCase struct {
	input string
}

type refResult struct {
	impossible bool
	length     int
}

type candResult struct {
	impossible bool
	length     int
	strings    []string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		refRes, err := parseReference(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, candOut)
			os.Exit(1)
		}
		candRes, err := parseCandidate(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot parse output on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, candOut)
			os.Exit(1)
		}

		if err := judge(tc.input, refRes, candRes); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n", i+1, err, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "193C-ref-*")
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

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseReference(out string) (refResult, error) {
	res := refResult{}
	trimmed := strings.TrimSpace(out)
	if trimmed == "-1" {
		res.impossible = true
		return res, nil
	}
	reader := bufio.NewReader(strings.NewReader(out))
	var length int
	if _, err := fmt.Fscan(reader, &length); err != nil {
		return res, err
	}
	if length < 0 {
		return res, fmt.Errorf("negative length %d", length)
	}
	res.length = length
	return res, nil
}

func parseCandidate(out string) (candResult, error) {
	res := candResult{}
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return res, fmt.Errorf("empty output")
	}
	if trimmed == "-1" {
		res.impossible = true
		return res, nil
	}
	reader := bufio.NewReader(strings.NewReader(out))
	var lengthStr string
	if _, err := fmt.Fscan(reader, &lengthStr); err != nil {
		return res, err
	}
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return res, fmt.Errorf("cannot parse length: %v", err)
	}
	if length < 0 {
		return res, fmt.Errorf("negative length %d", length)
	}
	res.length = length
	res.strings = make([]string, 4)
	for i := 0; i < 4; i++ {
		if _, err := fmt.Fscan(reader, &res.strings[i]); err != nil {
			return res, fmt.Errorf("missing string %d: %v", i+1, err)
		}
	}
	return res, nil
}

func judge(input string, refRes refResult, cand candResult) error {
	target, err := parseInput(input)
	if err != nil {
		return err
	}
	if refRes.impossible {
		if cand.impossible {
			return nil
		}
		return fmt.Errorf("expected -1 but got a construction")
	}
	if cand.impossible {
		return fmt.Errorf("candidate printed -1 but solution exists")
	}
	if cand.length != refRes.length {
		return fmt.Errorf("expected length %d, got %d", refRes.length, cand.length)
	}
	if len(cand.strings) != 4 {
		return fmt.Errorf("need 4 strings, got %d", len(cand.strings))
	}
	for i, s := range cand.strings {
		if len(s) != cand.length {
			return fmt.Errorf("string %d has length %d (expected %d)", i+1, len(s), cand.length)
		}
		for _, ch := range s {
			if ch != 'a' && ch != 'b' {
				return fmt.Errorf("string %d has invalid character %q", i+1, ch)
			}
		}
	}
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			dist := hamming(cand.strings[i], cand.strings[j])
			if dist != target[i][j] {
				return fmt.Errorf("h(s%d,s%d) = %d, expected %d", i+1, j+1, dist, target[i][j])
			}
		}
	}
	return nil
}

func parseInput(input string) ([4][4]int, error) {
	var d12, d13, d14, d23, d24, d34 int
	reader := strings.NewReader(input)
	if _, err := fmt.Fscan(reader, &d12, &d13, &d14, &d23, &d24, &d34); err != nil {
		return [4][4]int{}, err
	}
	var d [4][4]int
	d[0][1], d[1][0] = d12, d12
	d[0][2], d[2][0] = d13, d13
	d[0][3], d[3][0] = d14, d14
	d[1][2], d[2][1] = d23, d23
	d[1][3], d[3][1] = d24, d24
	d[2][3], d[3][2] = d34, d34
	return d, nil
}

func hamming(a, b string) int {
	cnt := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			cnt++
		}
	}
	return cnt
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(19301930))
	var tests []testCase

	tests = append(tests, fromStrings("a", "b", "a", "b"))
	tests = append(tests, fromStrings("ab", "ba", "aa", "bb"))
	tests = append(tests, fromStrings("aaaa", "bbbb", "abab", "baba"))

	tests = append(tests, customInput(0, 0, 1, 1, 1, 1))
	tests = append(tests, customInput(5, 1, 2, 1, 3, 4))

	tests = append(tests, customInput(0, 0, 2, 1, 2, 3))

	for i := 0; i < 25; i++ {
		length := rng.Intn(30) + 1
		tests = append(tests, randomValidTest(rng, length))
	}

	for i := 0; i < 10; i++ {
		length := rng.Intn(70) + 30
		tests = append(tests, randomValidTest(rng, length))
	}

	for i := 0; i < 10; i++ {
		tests = append(tests, randomNumbersTest(rng))
	}

	return tests
}

func fromStrings(s1, s2, s3, s4 string) testCase {
	strs := []string{s1, s2, s3, s4}
	var vals [6]int
	idx := 0
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			vals[idx] = hamming(strs[i], strs[j])
			idx++
		}
	}
	return testCase{input: formatInput(vals)}
}

func customInput(d12, d13, d14, d23, d24, d34 int) testCase {
	return testCase{input: formatInput([6]int{d12, d13, d14, d23, d24, d34})}
}

func formatInput(vals [6]int) string {
	return fmt.Sprintf("%d %d %d\n%d %d\n%d\n", vals[0], vals[1], vals[2], vals[3], vals[4], vals[5])
}

func randomValidTest(rng *rand.Rand, length int) testCase {
	strs := make([]string, 4)
	for i := 0; i < 4; i++ {
		b := make([]byte, length)
		for j := 0; j < length; j++ {
			if rng.Intn(2) == 0 {
				b[j] = 'a'
			} else {
				b[j] = 'b'
			}
		}
		strs[i] = string(b)
	}
	allEqual := true
outer:
	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			if strs[i] != strs[j] {
				allEqual = false
				break outer
			}
		}
	}
	if allEqual {
		b := []byte(strs[1])
		if b[0] == 'a' {
			b[0] = 'b'
		} else {
			b[0] = 'a'
		}
		strs[1] = string(b)
	}
	return fromStrings(strs[0], strs[1], strs[2], strs[3])
}

func randomNumbersTest(rng *rand.Rand) testCase {
	var vals [6]int
	for i := 0; i < 6; i++ {
		vals[i] = rng.Intn(60)
	}
	zero := true
	for _, v := range vals {
		if v != 0 {
			zero = false
			break
		}
	}
	if zero {
		vals[0] = 1
	}
	return testCase{input: formatInput(vals)}
}
