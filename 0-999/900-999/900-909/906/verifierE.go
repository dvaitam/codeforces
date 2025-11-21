package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type segment struct {
	l, r int
}

type solution struct {
	impossible bool
	k          int
	segments   []segment
}

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refSrc, err := locateReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refBin, err := buildReference(refSrc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refSol, err := parseSolution(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\nOutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		candSol, err := parseSolution(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse candidate output on test %d: %v\nOutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		s, t, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid test input on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if err := validateSolutions(refSol, candSol, s, t); err != nil {
			fmt.Fprintf(os.Stderr, "candidate invalid on test %d: %v\nInput:\n%s\nCandidate output:\n%s\n", idx+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"906E.go",
		filepath.Join("0-999", "900-999", "900-909", "906", "906E.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 906E.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref906E_%d.bin", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", outPath, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return outPath, nil
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
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseSolution(out string) (solution, error) {
	var sol solution
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return sol, fmt.Errorf("empty output")
	}
	reader := strings.NewReader(out)
	var first string
	if _, err := fmt.Fscan(reader, &first); err != nil {
		return sol, fmt.Errorf("failed to read first token: %v", err)
	}
	if first == "-1" {
		sol.impossible = true
		return sol, nil
	}
	k, err := strconv.Atoi(first)
	if err != nil {
		return sol, fmt.Errorf("invalid k: %v", err)
	}
	if k < 0 {
		return sol, fmt.Errorf("negative k")
	}
	sol.k = k
	sol.segments = make([]segment, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &sol.segments[i].l, &sol.segments[i].r); err != nil {
			return sol, fmt.Errorf("failed to read segment %d: %v", i+1, err)
		}
	}
	return sol, nil
}

func parseInput(input string) (string, string, error) {
	reader := strings.NewReader(input)
	var s, t string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return "", "", fmt.Errorf("failed to read s: %v", err)
	}
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return "", "", fmt.Errorf("failed to read t: %v", err)
	}
	return s, t, nil
}

func validateSolutions(refSol, candSol solution, s, t string) error {
	if refSol.impossible {
		if candSol.impossible {
			return nil
		}
		return fmt.Errorf("reference determined impossibility but candidate provided a solution")
	}
	if candSol.impossible {
		return fmt.Errorf("candidate claims impossibility but reference found a solution")
	}
	if candSol.k != refSol.k {
		return fmt.Errorf("non-minimal number of substrings: expected %d, got %d", refSol.k, candSol.k)
	}
	if err := validateSegments(candSol.segments, len(s)); err != nil {
		return err
	}
	res := applySegments(t, candSol.segments)
	if res != s {
		return fmt.Errorf("applied substrings do not transform t into s")
	}
	return nil
}

func validateSegments(segs []segment, n int) error {
	for idx, seg := range segs {
		if seg.l < 1 || seg.r > n || seg.l > seg.r {
			return fmt.Errorf("segment %d is out of range: [%d, %d]", idx+1, seg.l, seg.r)
		}
	}
	sorted := append([]segment(nil), segs...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].l == sorted[j].l {
			return sorted[i].r < sorted[j].r
		}
		return sorted[i].l < sorted[j].l
	})
	prevEnd := 0
	for _, seg := range sorted {
		if seg.l <= prevEnd {
			return fmt.Errorf("segments [%d, %d] and previous overlap", seg.l, seg.r)
		}
		prevEnd = seg.r
	}
	return nil
}

func applySegments(t string, segs []segment) string {
	data := []byte(t)
	for _, seg := range segs {
		l := seg.l - 1
		r := seg.r - 1
		for l < r {
			data[l], data[r] = data[r], data[l]
			l++
			r--
		}
	}
	return string(data)
}

func generateTests() []testCase {
	tests := []testCase{
		newTest("a", "a"),
		newTest("ab", "ba"),
		newTest("ab", "ac"),
		newTest("abcd", "abcd"),
		newTest("abcd", "badc"),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 50 {
		tests = append(tests, randomConvertibleTest(rng, rng.Intn(20)+1))
	}
	for len(tests) < 90 {
		tests = append(tests, randomConvertibleTest(rng, rng.Intn(200)+50))
	}
	tests = append(tests, randomConvertibleTest(rng, 1000))
	tests = append(tests, randomConvertibleTest(rng, 5000))
	tests = append(tests, randomConvertibleTest(rng, 50000))
	tests = append(tests, randomConvertibleTest(rng, 200000))
	tests = append(tests, randomImpossibleTest(rng, 200))
	tests = append(tests, randomImpossibleTest(rng, 1000))
	return tests
}

func randomConvertibleTest(rng *rand.Rand, n int) testCase {
	if n < 1 {
		n = 1
	}
	s := randomString(rng, n)
	tBytes := []byte(s)
	if n >= 2 {
		maxSeg := determineMaxSegments(n)
		segCount := 0
		if maxSeg > 0 {
			segCount = rng.Intn(maxSeg + 1)
		}
		segments := generateNonOverlappingSegments(rng, n, segCount)
		for _, seg := range segments {
			l := seg.l - 1
			r := seg.r - 1
			for l < r {
				tBytes[l], tBytes[r] = tBytes[r], tBytes[l]
				l++
				r--
			}
		}
	}
	return newTestFromStrings(s, string(tBytes))
}

func determineMaxSegments(n int) int {
	switch {
	case n >= 10000:
		return 10
	case n >= 1000:
		return 7
	case n >= 200:
		return 5
	case n >= 50:
		return 3
	case n >= 2:
		return 2
	default:
		return 0
	}
}

func generateNonOverlappingSegments(rng *rand.Rand, n, count int) []segment {
	var segments []segment
	if count <= 0 {
		return segments
	}
	used := make([]bool, n)
	attempts := 0
	maxAttempts := count*20 + 20
	for len(segments) < count && attempts < maxAttempts {
		attempts++
		l := rng.Intn(n)
		if l+1 >= n {
			continue
		}
		maxLen := n - l
		if maxLen < 2 {
			continue
		}
		length := rng.Intn(maxLen-1) + 2
		r := l + length - 1
		if hasOverlap(used, l, r) {
			continue
		}
		for i := l; i <= r; i++ {
			used[i] = true
		}
		segments = append(segments, segment{l + 1, r + 1})
	}
	return segments
}

func hasOverlap(used []bool, l, r int) bool {
	for i := l; i <= r; i++ {
		if used[i] {
			return true
		}
	}
	return false
}

func randomImpossibleTest(rng *rand.Rand, n int) testCase {
	if n < 1 {
		n = 1
	}
	s := randomString(rng, n)
	tBytes := []byte(s)
	pos := rng.Intn(n)
	orig := tBytes[pos]
	for {
		tBytes[pos] = byte('a' + rng.Intn(26))
		if tBytes[pos] != orig {
			break
		}
	}
	return newTestFromStrings(s, string(tBytes))
}

func randomString(rng *rand.Rand, n int) string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(alphabet[rng.Intn(len(alphabet))])
	}
	return b.String()
}

func newTestFromStrings(s, t string) testCase {
	return testCase{input: fmt.Sprintf("%s\n%s\n", s, t)}
}

func newTest(s, t string) testCase {
	return newTestFromStrings(s, t)
}
