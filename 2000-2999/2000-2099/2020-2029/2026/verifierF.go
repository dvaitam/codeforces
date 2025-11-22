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
	input    string
	outCount int
}

type ice struct {
	price int
	tasty int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		expectOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		expectVals, err := parseOutputs(expectOut, tc.outCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, expectOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, tc.outCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nInput:\n%sOutput:\n%s", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}

		if len(expectVals) != len(gotVals) {
			fmt.Fprintf(os.Stderr, "test %d: output count mismatch, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, len(expectVals), len(gotVals), tc.input, gotOut)
			os.Exit(1)
		}
		for i := range expectVals {
			if expectVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "test %d: mismatch at answer %d, expected %d got %d\nInput:\n%sOutput:\n%s", idx+1, i+1, expectVals[i], gotVals[i], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"2026F.go",
		filepath.Join("2000-2999", "2000-2099", "2020-2029", "2026", "2026F.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 2026F.go")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref2026F_%d.bin", time.Now().UnixNano()))
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

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = val
	}
	return res, nil
}

func generateTests() []testCase {
	tests := make([]testCase, 0, 40)
	tests = append(tests, sampleLikeTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 25; i++ {
		tests = append(tests, randomTest(rng))
	}
	return tests
}

func sampleLikeTests() []testCase {
	var res []testCase

	// Simple add/query/clone flow.
	queries1 := []string{
		"2 1 5 7",
		"2 1 3 4",
		"4 1 4",
		"4 1 8",
		"1 1",
		"4 2 4",
		"4 2 8",
	}
	res = append(res, buildTest(queries1))

	// Deletions and cloning after deletions.
	queries2 := []string{
		"2 1 5 10",
		"2 1 3 3",
		"3 1",
		"4 1 10",
		"1 1",
		"4 2 10",
	}
	res = append(res, buildTest(queries2))

	// Query on empty store.
	res = append(res, buildTest([]string{"4 1 2000"}))

	return res
}

func buildTest(queries []string) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	outCnt := 0
	for _, line := range queries {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == "4" {
			outCnt++
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), outCount: outCnt}
}

func randomTest(rng *rand.Rand) testCase {
	q := rng.Intn(250) + 50    // between 50 and 299 queries
	stores := make([][]ice, 2) // 1-based; index 1 exists
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	outCnt := 0

	for i := 0; i < q; i++ {
		// Guarantee at least one type 4.
		if i == q-1 && outCnt == 0 {
			sb.WriteString(fmt.Sprintf("4 1 %d\n", rng.Intn(2000)+1))
			outCnt++
			continue
		}

		tp := rng.Intn(100)
		switch {
		case tp < 20: // type 1
			x := rng.Intn(len(stores)-1) + 1
			stores = append(stores, append([]ice(nil), stores[x]...))
			sb.WriteString(fmt.Sprintf("1 %d\n", x))
		case tp < 55: // type 2
			x := rng.Intn(len(stores)-1) + 1
			p := rng.Intn(2000) + 1
			t := rng.Intn(2000) + 1
			stores[x] = append(stores[x], ice{price: p, tasty: t})
			sb.WriteString(fmt.Sprintf("2 %d %d %d\n", x, p, t))
		case tp < 70: // type 3, only if store not empty
			x := rng.Intn(len(stores)-1) + 1
			if len(stores[x]) == 0 {
				// fallback to type2
				p := rng.Intn(2000) + 1
				t := rng.Intn(2000) + 1
				stores[x] = append(stores[x], ice{price: p, tasty: t})
				sb.WriteString(fmt.Sprintf("2 %d %d %d\n", x, p, t))
			} else {
				stores[x] = stores[x][1:]
				sb.WriteString(fmt.Sprintf("3 %d\n", x))
			}
		default: // type 4
			x := rng.Intn(len(stores)-1) + 1
			p := rng.Intn(2000) + 1
			sb.WriteString(fmt.Sprintf("4 %d %d\n", x, p))
			outCnt++
		}
	}

	return testCase{input: sb.String(), outCount: outCnt}
}
