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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, wantOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if len(got) != len(want) {
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
		"2149D.go",
		filepath.Join("2000-2999", "2100-2199", "2140-2149", "2149", "2149D.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2149D.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2149D_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) ([]int64, error) {
	fields := strings.Fields(out)
	res := make([]int64, len(fields))
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
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
		buildTest([]string{"a"}, []int{1}),
		buildTest([]string{"ab", "ba"}, []int{2, 2}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 50 {
		t := rng.Intn(10) + 1
		tStr := make([]string, t)
		nArr := make([]int, t)
		for i := 0; i < t; i++ {
			n := rng.Intn(20) + 1
			nArr[i] = n
			var sb strings.Builder
			for j := 0; j < n; j++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('a')
				} else {
					sb.WriteByte('b')
				}
			}
			tStr[i] = sb.String()
		}
		tests = append(tests, buildTest(tStr, nArr))
	}
	return tests
}

func buildTest(strs []string, lengths []int) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(strs)))
	for i := range strs {
		b.WriteString(fmt.Sprintf("%d\n%s\n", lengths[i], strs[i]))
	}
	return testCase{input: b.String()}
}
