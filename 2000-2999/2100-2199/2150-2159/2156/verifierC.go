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

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
		wantOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		want, err := parseOutput(wantOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if len(want) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d answers, got %d\nInput:\n%s\n", idx+1, len(want), len(got), tc.input)
			os.Exit(1)
		}
		for i := range want {
			if want[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d mismatch at case %d: expected %d, got %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, i+1, want[i], got[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2156C.go",
		filepath.Join("2000-2999", "2100-2199", "2150-2159", "2156", "2156C.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2156C.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2156C_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) ([]int, error) {
	fields := strings.Fields(out)
	res := make([]int, len(fields))
	for i, tok := range fields {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([]int{1}, 0, [][]int{{1}}),
		buildTest([]int{3, 1, 2}, 1, [][]int{{1, 3}, {2, 3}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(15) + 1
		k := rng.Intn(n)
		tests = append(tests, buildTest(randArray(rng, n), k, nil))
	}
	return tests
}

func randArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	return arr
}

func buildTest(arr []int, k int, extra [][]int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("1\n%d %d\n", len(arr), k))
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(strconv.Itoa(v))
	}
	b.WriteByte('\n')
	for _, row := range extra {
		for i, v := range row {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(v))
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}
