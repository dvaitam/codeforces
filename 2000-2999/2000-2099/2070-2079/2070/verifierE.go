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

const refSource = "2070E.go"

type testCase struct {
	name string
	s    string
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candPath := os.Args[len(os.Args)-1]

	refBin, cleanupRef, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := prepareCandidate(candPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to prepare candidate:", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n%s\n", len(tc.s), tc.s)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput preview: %s\n", idx+1, tc.name, err, summarizeString(tc.s))
			os.Exit(1)
		}
		refAns, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput preview: %s\noutput:\n%s\n", idx+1, tc.name, err, summarizeString(tc.s), candOut)
			os.Exit(1)
		}
		candAns, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput preview: %s\noutput:\n%s\n", idx+1, tc.name, err, summarizeString(tc.s), candOut)
			os.Exit(1)
		}

		if candAns != refAns {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d got %d\ninput (n=%d): %s\n", idx+1, tc.name, refAns, candAns, len(tc.s), summarizeString(tc.s))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)

	tmpDir, err := os.MkdirTemp("", "ref2070E-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")

	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func prepareCandidate(path string) (string, func(), error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", nil, err
	}
	if filepath.Ext(abs) != ".go" {
		return abs, func() {}, nil
	}

	tmpDir, err := os.MkdirTemp("", "cand2070E-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "candidate")

	cmd := exec.Command("go", "build", "-o", binPath, abs)
	cmd.Dir = filepath.Dir(abs)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("candidate build failed: %v\n%s", err, stderr.String())
	}

	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string) (int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != 1 {
		return 0, fmt.Errorf("expected 1 integer, got %d", len(tokens))
	}
	v, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", tokens[0], err)
	}
	return v, nil
}

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []testCase{
		{name: "sample", s: "0010010011"},
	}
	tests = append(tests, enumerateSmall()...)
	tests = append(tests, manualPatterns()...)
	tests = append(tests, randomTests(rng, 12, 5, 20, "rand_small")...)
	tests = append(tests, randomTests(rng, 12, 30, 200, "rand_mid")...)
	tests = append(tests, randomTests(rng, 8, 500, 1200, "rand_big")...)
	tests = append(tests, stressTests(rng)...)
	return tests
}

func enumerateSmall() []testCase {
	var res []testCase
	for n := 1; n <= 4; n++ {
		limit := 1 << n
		for mask := 0; mask < limit; mask++ {
			var sb strings.Builder
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
			res = append(res, testCase{
				name: fmt.Sprintf("len%d_mask%d", n, mask),
				s:    sb.String(),
			})
		}
	}
	return res
}

func manualPatterns() []testCase {
	return []testCase{
		{name: "all_zero_6", s: strings.Repeat("0", 6)},
		{name: "all_one_6", s: strings.Repeat("1", 6)},
		{name: "alternating_10", s: "0101010101"},
		{name: "alternating_shifted_11", s: "10101010101"},
		{name: "blocks_growth", s: "0000111100001111"},
		{name: "dense_ones_tail", s: "0000011111"},
		{name: "dense_zeros_tail", s: "1111100000"},
		{name: "two_large_blocks", s: strings.Repeat("0", 30) + strings.Repeat("1", 30)},
		{name: "pivoted_block", s: "11110000111100001111"},
	}
}

func randomTests(rng *rand.Rand, count, minLen, maxLen int, prefix string) []testCase {
	res := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxLen-minLen+1) + minLen
		res = append(res, testCase{
			name: fmt.Sprintf("%s_%d_len_%d", prefix, i+1, n),
			s:    randBinary(rng, n),
		})
	}
	return res
}

func stressTests(rng *rand.Rand) []testCase {
	return []testCase{
		{name: "all_zero_20000", s: strings.Repeat("0", 20000)},
		{name: "all_one_20000", s: strings.Repeat("1", 20000)},
		{name: "long_blocks", s: buildBlocks(50000, 7)},
		{name: "random_100k", s: randBinary(rng, 100000)},
		{name: "pattern_150k", s: buildBlocks(150000, 3)},
		{name: "random_300k", s: randBinary(rng, 300000)},
	}
}

func randBinary(rng *rand.Rand, n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func buildBlocks(total, block int) string {
	var sb strings.Builder
	sb.Grow(total)
	cur := byte('0')
	remaining := total
	for remaining > 0 {
		lenBlock := block
		if lenBlock > remaining {
			lenBlock = remaining
		}
		for i := 0; i < lenBlock; i++ {
			sb.WriteByte(cur)
		}
		remaining -= lenBlock
		if cur == '0' {
			cur = '1'
		} else {
			cur = '0'
		}
		block++
	}
	return sb.String()
}

func summarizeString(s string) string {
	if len(s) <= 80 {
		return s
	}
	return fmt.Sprintf("%s...%s (len=%d)", s[:30], s[len(s)-30:], len(s))
}
