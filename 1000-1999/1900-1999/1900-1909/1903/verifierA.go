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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
				fmt.Fprintf(os.Stderr, "test %d mismatch on case %d: expected %s, got %s\nInput:\n%s\nCandidate output:\n%s\n", idx+1, i+1, want[i], got[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1903A.go",
		filepath.Join("1000-1999", "1900-1999", "1900-1909", "1903", "1903A.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1903A.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1903A_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	for i, f := range fields {
		fields[i] = strings.ToUpper(f)
		if fields[i] != "YES" && fields[i] != "NO" {
			return nil, fmt.Errorf("invalid token %q", f)
		}
	}
	return fields, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([]testInput{
			{n: 1, k: 1, arr: []int{5}},
			{n: 2, k: 1, arr: []int{2, 1}},
		}),
		buildTest([]testInput{
			{n: 3, k: 2, arr: []int{3, 1, 2}},
			{n: 5, k: 1, arr: []int{1, 2, 3, 4, 5}},
			{n: 4, k: 1, arr: []int{4, 1, 2, 3}},
		}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 60 {
		t := rng.Intn(20) + 1
		var cases []testInput
		for i := 0; i < t; i++ {
			n := rng.Intn(10) + 1
			k := rng.Intn(5) + 1
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rng.Intn(100)
			}
			cases = append(cases, testInput{n: n, k: k, arr: arr})
		}
		tests = append(tests, buildTest(cases))
	}
	return tests
}

type testInput struct {
	n   int
	k   int
	arr []int
}

func buildTest(cases []testInput) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		b.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.k))
		for i, v := range cs.arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(v))
		}
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}
