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

const referenceSource = "2000-2999/2000-2099/2060-2069/2061/2061D.go"

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	candidatePath, err := candidatePathFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	candidate, candCleanup, err := prepareCandidateBinary(candidatePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	refPath := referencePath()
	refBin, refCleanup, err := buildReferenceBinary(refPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		expect, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\nraw:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}
		have, err := parseAnswers(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\nraw:\n%s\n", idx+1, tc.name, err, candOut)
			fmt.Fprintln(os.Stderr, previewInput(tc.input))
			os.Exit(1)
		}

		for i := 0; i < tc.t; i++ {
			if expect[i] != have[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %s got %s\n", idx+1, tc.name, i+1, expect[i], have[i])
				fmt.Fprintln(os.Stderr, previewInput(tc.input))
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func candidatePathFromArgs() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("usage: go run verifierD.go /path/to/binary-or-source")
	}
	return os.Args[1], nil
}

func referencePath() string {
	if _, file, _, ok := runtime.Caller(0); ok {
		return filepath.Join(filepath.Dir(file), "2061D.go")
	}
	return referenceSource
}

func buildReferenceBinary(src string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "verifier2061D")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "ref2061d")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func prepareCandidateBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmpDir, err := os.MkdirTemp("", "candidate2061D")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	bin := filepath.Join(tmpDir, "candidate")
	cmd := exec.Command("go", "build", "-o", bin, path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build candidate: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseAnswers(out string, t int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]string, t)
	for i, tok := range tokens {
		l := strings.ToLower(tok)
		if l == "yes" {
			res[i] = "yes"
		} else if l == "no" {
			res[i] = "no"
		} else {
			return nil, fmt.Errorf("token %q is not yes/no", tok)
		}
	}
	return res, nil
}

func buildTests() []testCase {
	var tests []testCase

	tests = append(tests, testCase{
		name: "handmade-basic",
		t:    5,
		input: "5\n" +
			"2 1\n4 5\n9\n" +
			"2 1\n3 6\n9\n" +
			"4 2\n1 2 2 2\n3 4\n" +
			"4 2\n1 1 3 3\n5 5\n" +
			"5 5\n2 3 4 5 5\n2 3 4 5 5\n",
	})

	tests = append(tests, testCase{
		name: "small-handmade",
		t:    4,
		input: "4\n" +
			"1 1\n5\n5\n" +
			"3 3\n1 2 3\n1 2 3\n" +
			"3 1\n2 2 2\n6\n" +
			"5 2\n1 1 1 1 1\n2 3\n",
	})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 150; i++ {
		tests = append(tests, randomTest(rng, i, 40, 80))
	}
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTest(rng, i, 500, 50))
	}

	return tests
}

func randomTest(rng *rand.Rand, idx, maxN, maxT int) testCase {
	t := rng.Intn(maxT-1) + 2
	if t > 200 {
		t = 200
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(maxN-1) + 2
		m := rng.Intn(n) + 1
		if rng.Intn(5) == 0 {
			m = n
		}
		if rng.Intn(10) == 0 {
			m = 1
		}
		if rng.Intn(7) == 0 {
			n = 1
			m = 1
		}
		if m > n {
			m = n
		}
		a, b := makeTestCase(rng, n, m)
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for j, v := range b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return testCase{
		name:  fmt.Sprintf("random-%d-%d", maxN, idx+1),
		t:     t,
		input: sb.String(),
	}
}

func makeTestCase(rng *rand.Rand, n, m int) ([]int, []int) {
	if n == 1 && m == 1 {
		val := rng.Intn(10) + 1
		return []int{val}, []int{val}
	}

	if rng.Intn(2) == 0 {
		// Guaranteed YES: keep all numbers within base..base+1 so merges are always allowed.
		base := rng.Intn(1_000_000_000-2) + 1
		a := make([]int, n)
		for i := range a {
			if rng.Intn(2) == 0 {
				a[i] = base
			} else {
				a[i] = base + 1
			}
		}
		cur := append([]int(nil), a...)
		for len(cur) > m {
			i := rng.Intn(len(cur))
			j := rng.Intn(len(cur) - 1)
			if j >= i {
				j++
			}
			if abs(cur[i]-cur[j]) > 1 {
				cur[i] = base
				cur[j] = base
			}
			sum := cur[i] + cur[j]
			if i > j {
				i, j = j, i
			}
			cur[i] = sum
			cur = append(cur[:j], cur[j+1:]...)
		}
		return a, cur
	}

	// Mixed/likely NO cases.
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(1_000_000_000) + 1
	}
	b := make([]int, m)
	for i := range b {
		b[i] = rng.Intn(1_000_000_000) + 1
	}
	if rng.Intn(3) == 0 {
		// Force sum mismatch to be definite NO.
		sumA := 0
		for _, v := range a {
			sumA += v
		}
		sumB := 0
		for _, v := range b {
			sumB += v
		}
		if sumA == sumB {
			b[0]++
		}
	}
	return a, b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func previewInput(in string) string {
	const limit = 500
	if len(in) <= limit {
		return in
	}
	return in[:limit] + "..."
}
