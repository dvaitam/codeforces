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
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/binary")
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
		if want != got {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, want, got, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"1662L.go",
		filepath.Join("1000-1999", "1600-1699", "1660-1669", "1662", "1662L.go"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not find 1662L.go reference")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref1662L_%d.bin", time.Now().UnixNano()))
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

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest(1, 1, []int{0}, []int64{0}),
		buildTest(2, 1, []int{1, 2}, []int64{1, 2}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		n := rng.Intn(50) + 1
		tests = append(tests, randomTest(rng, n))
	}
	tests = append(tests, randomTest(rng, 2000))
	return tests
}

func buildTest(n int, v int64, tArr []int, aArr []int64) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d %d\n", n, v))
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i < len(tArr) {
			b.WriteString(strconv.Itoa(tArr[i]))
		} else {
			b.WriteString("0")
		}
	}
	b.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		if i < len(aArr) {
			b.WriteString(fmt.Sprintf("%d", aArr[i]))
		} else {
			b.WriteString("0")
		}
	}
	b.WriteByte('\n')
	return testCase{input: b.String()}
}

func randomTest(rng *rand.Rand, n int) testCase {
	v := rng.Int63n(1_000_000_000) - 500_000_000
	if v == 0 {
		v = 1
	}
	tArr := make([]int, n)
	aArr := make([]int64, n)
	for i := 0; i < n; i++ {
		tArr[i] = rng.Intn(1_000_000_000)
		aArr[i] = rng.Int63n(1_000_000_000) - 500_000_000
	}
	return buildTest(n, v, tArr, aArr)
}
