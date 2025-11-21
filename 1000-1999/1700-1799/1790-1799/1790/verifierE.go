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
	"time"
)

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

		if len(got) != len(want) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d answers, got %d\nInput:\n%s\n", idx+1, len(want), len(got), tc.input)
			os.Exit(1)
		}
		if err := compare(tc, want, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\nInput:\n%s\nCandidate output:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1790E.go",
		filepath.Join("1000-1999", "1700-1799", "1790-1799", "1790", "1790E.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1790E.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1790E_%d.bin", time.Now().UnixNano()))
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

type answer struct {
	valid bool
	a, b  int64
}

func parseOutput(out string) ([]answer, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	var res []answer
	for sc.Scan() {
		token := sc.Text()
		if token == "-1" {
			res = append(res, answer{})
			continue
		}
		a, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("missing second integer for pair")
		}
		bval, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", sc.Text())
		}
		res = append(res, answer{valid: true, a: a, b: bval})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}
	return res, nil
}

func compare(tc testCase, want, got []answer) error {
	sc := bufio.NewScanner(strings.NewReader(tc.input))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return fmt.Errorf("invalid test input")
	}
	t, _ := strconv.Atoi(sc.Text())
	xVals := make([]int64, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return fmt.Errorf("invalid test input")
		}
		x, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid x %q", sc.Text())
		}
		xVals[i] = x
	}

	for i := 0; i < t; i++ {
		x := xVals[i]
		w := want[i]
		c := got[i]
		if !w.valid {
			if c.valid {
				return fmt.Errorf("query %d: candidate claims solution but reference says -1", i+1)
			}
			continue
		}
		if !c.valid {
			return fmt.Errorf("query %d: candidate output -1 but solution exists", i+1)
		}
		if c.a^c.b != x {
			return fmt.Errorf("query %d: a xor b != x", i+1)
		}
		if c.a+c.b != 2*x {
			return fmt.Errorf("query %d: a + b != 2*x", i+1)
		}
		if c.a < 0 || c.b < 0 {
			return fmt.Errorf("query %d: negative values not allowed", i+1)
		}
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([]int64{0}),
		buildTest([]int64{1}),
		buildTest([]int64{2, 4, 6}),
		buildTest([]int64{7, 8, 1024}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		t := rng.Intn(10) + 1
		arr := make([]int64, t)
		for i := 0; i < t; i++ {
			arr[i] = rng.Int63n(1_000_000_000)
		}
		tests = append(tests, buildTest(arr))
	}
	return tests
}

func buildTest(xs []int64) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(xs)))
	for _, x := range xs {
		b.WriteString(fmt.Sprintf("%d\n", x))
	}
	return testCase{input: b.String()}
}
