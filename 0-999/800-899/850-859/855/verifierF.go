package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const maxX = 100000

type testCase struct {
	input string
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
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, tc.input, err)
			os.Exit(1)
		}
		if normalizeOutput(got) != normalizeOutput(want) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nInput:\n%sExpected:\n%sGot:\n%s\n", idx+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func locateReference() (string, error) {
	candidates := []string{
		"855F.go",
		filepath.Join("0-999", "800-899", "850-859", "855", "855F.go"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("could not find 855F.go relative to working directory")
}

func buildReference(src string) (string, error) {
	outPath := filepath.Join(os.TempDir(), fmt.Sprintf("ref855F_%d.bin", time.Now().UnixNano()))
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

func normalizeOutput(out string) string {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		newTest([]string{
			"1 1 2 5",
			"2 1 3",
		}),
		newTest([]string{
			"1 1 5 10",
			"1 1 5 -7",
			"2 1 5",
		}),
		newTest([]string{
			"1 1 4 3",
			"1 2 5 -2",
			"1 2 5 6",
			"2 1 6",
			"2 2 4",
		}),
	)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 80 {
		tests = append(tests, randomTestCase(rng, rng.Intn(200)+1))
	}
	tests = append(tests, randomTestCase(rng, 50000))
	return tests
}

func randomTestCase(rng *rand.Rand, q int) testCase {
	if q < 2 {
		q = 2
	}
	ops := make([]string, q)
	hasQuery := false
	for i := 0; i < q; i++ {
		if rng.Intn(3) == 0 || (!hasQuery && i == q-1) {
			l := rng.Intn(maxX-1) + 1
			r := l + rng.Intn(maxX-l) + 1
			k := rng.Int63n(2_000_000_000) - 1_000_000_000
			if k == 0 {
				if rng.Intn(2) == 0 {
					k = 1
				} else {
					k = -1
				}
			}
			ops[i] = fmt.Sprintf("1 %d %d %d", l, r, k)
		} else {
			l := rng.Intn(maxX-1) + 1
			r := l + rng.Intn(maxX-l) + 1
			ops[i] = fmt.Sprintf("2 %d %d", l, r)
			hasQuery = true
		}
	}
	return newTest(ops)
}

func newTest(ops []string) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(ops)))
	for _, op := range ops {
		b.WriteString(op)
		b.WriteByte('\n')
	}
	return testCase{input: b.String()}
}
