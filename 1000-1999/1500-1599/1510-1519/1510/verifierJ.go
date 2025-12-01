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

const refSource = "./1510J.go"

type testCase struct {
	input string
	mask  string
}

type result struct {
	possible bool
	profile  []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		refRes, err := parseResult(refOut)
		if err != nil {
			fail("could not parse reference output on test %d: %v", idx+1, err)
		}

		candOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fail("candidate crashed on test %d: %v", idx+1, err)
		}
		candRes, err := parseResult(candOut)
		if err != nil {
			fail("failed to parse candidate output on test %d: %v", idx+1, err)
		}

		if refRes.possible {
			if !candRes.possible {
				fail("candidate claims -1 on solvable test %d", idx+1)
			}
			if err := validateProfile(tc.mask, candRes.profile); err != nil {
				fail("invalid profile on test %d: %v", idx+1, err)
			}
		} else {
			if candRes.possible {
				fail("candidate outputs a profile on impossible test %d", idx+1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1510J-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, buf.String())
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runCmd(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runCmd(cmd, input)
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

func runCmd(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return stdout.String(), fmt.Errorf("%v\nstderr:\n%s", err, stderr.String())
		}
		return stdout.String(), err
	}
	return stdout.String(), nil
}

func parseResult(out string) (result, error) {
	reader := strings.NewReader(out)
	var first string
	if _, err := fmt.Fscan(reader, &first); err != nil {
		return result{}, fmt.Errorf("missing verdict: %w", err)
	}
	if first == "-1" {
		return result{possible: false}, nil
	}
	k, err := strconv.Atoi(first)
	if err != nil {
		return result{}, fmt.Errorf("invalid k value: %w", err)
	}
	if k < 0 {
		return result{}, fmt.Errorf("negative k")
	}
	profile := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &profile[i]); err != nil {
			return result{}, fmt.Errorf("failed to read profile entry %d: %w", i+1, err)
		}
		if profile[i] <= 0 {
			return result{}, fmt.Errorf("profile entries must be positive")
		}
	}
	return result{possible: true, profile: profile}, nil
}

func validateProfile(mask string, profile []int) error {
	built, ok := buildMask(len(mask), profile)
	if !ok {
		return fmt.Errorf("profile cannot produce any valid mask of length %d", len(mask))
	}
	if built != mask {
		return fmt.Errorf("mask mismatch: expected %s, got %s", mask, built)
	}
	return nil
}

func buildMask(n int, profile []int) (string, bool) {
	if n < 0 {
		return "", false
	}
	if len(profile) == 0 {
		return strings.Repeat("_", n), true
	}

	earliest := make([]int, len(profile))
	pos := 0
	for i, v := range profile {
		earliest[i] = pos
		pos += v
		if pos > n {
			return "", false
		}
		if i+1 < len(profile) {
			pos++
		}
	}
	if pos > n {
		return "", false
	}

	latest := make([]int, len(profile))
	pos = n
	for i := len(profile) - 1; i >= 0; i-- {
		pos -= profile[i]
		if pos < 0 {
			return "", false
		}
		latest[i] = pos
		if i > 0 {
			pos--
			if pos < 0 {
				return "", false
			}
		}
	}

	res := make([]byte, n)
	for i := range res {
		res[i] = '_'
	}

	for i := 0; i < len(profile); i++ {
		if earliest[i] > latest[i] {
			return "", false
		}
		diff := latest[i] - earliest[i]
		forced := profile[i] - diff
		if forced <= 0 {
			continue
		}
		start := latest[i]
		end := start + forced
		if end > n {
			return "", false
		}
		for j := start; j < end; j++ {
			res[j] = '#'
		}
	}
	return string(res), true
}

func buildTests() []testCase {
	base := []string{
		"__#_____",
		"_#",
		"___",
		"#",
		"_",
		"#__#",
		"#_#_#",
	}
	tests := make([]testCase, 0, 32)
	for _, m := range base {
		tests = append(tests, testCase{input: m + "\n", mask: m})
	}

	rng := rand.New(rand.NewSource(1))
	for len(tests) < 20 {
		mask := randomProfileMask(rng)
		tests = append(tests, testCase{input: mask + "\n", mask: mask})
	}
	for len(tests) < 28 {
		mask := randomRawMask(rng)
		tests = append(tests, testCase{input: mask + "\n", mask: mask})
	}
	return tests
}

func randomProfileMask(rng *rand.Rand) string {
	n := rng.Intn(30) + 1
	maxBlocks := (n + 1) / 2
	k := rng.Intn(maxBlocks + 1)
	profile := make([]int, k)
	if k > 0 {
		remaining := n - (k - 1)
		for i := 0; i < k; i++ {
			minVal := 1
			maxVal := remaining - (k - 1 - i)
			if maxVal < minVal {
				maxVal = minVal
			}
			val := minVal
			if maxVal > minVal {
				val += rng.Intn(maxVal - minVal + 1)
			}
			profile[i] = val
			remaining -= val
		}
	}
	mask, ok := buildMask(n, profile)
	if !ok {
		panic("failed to build mask from generated profile")
	}
	return mask
}

func randomRawMask(rng *rand.Rand) string {
	n := rng.Intn(25) + 1
	var b strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b.WriteByte('#')
		} else {
			b.WriteByte('_')
		}
	}
	return b.String()
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
