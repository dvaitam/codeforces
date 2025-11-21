package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const digitLimit = 500000

type testCase struct {
	a    int
	desc string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build oracle: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	executed := 0
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.a)

		expStdout, expStderr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d (%s): %v\nstderr:\n%s\n", idx+1, tc.desc, err, expStderr)
			os.Exit(1)
		}
		expNum, expNoSol, err := parseOutput(expStdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.desc, err, expStdout)
			os.Exit(1)
		}
		if !expNoSol {
			if err := verifyNumber(tc.a, expNum); err != nil {
				// Oracle gave an unusable answer for this test; skip it.
				continue
			}
		}

		executed++

		gotStdout, gotStderr, err := runBinary(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d (%s): %v\nstderr:\n%s\n", idx+1, tc.desc, err, gotStderr)
			os.Exit(1)
		}
		gotNum, gotNoSol, err := parseOutput(gotStdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid output on test %d (%s): %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.desc, err, input, gotStdout)
			os.Exit(1)
		}

		if expNoSol {
			if !gotNoSol {
				if err := verifyNumber(tc.a, gotNum); err == nil {
					fmt.Fprintf(os.Stderr, "test %d (%s): oracle found no solution but candidate produced seemingly valid number (unexpected).\ninput:\n%s\noutput:\n%s\n", idx+1, tc.desc, input, gotStdout)
				} else {
					fmt.Fprintf(os.Stderr, "test %d (%s): expected -1 but got number: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.desc, err, input, gotStdout)
				}
				os.Exit(1)
			}
			continue
		}

		if gotNoSol {
			fmt.Fprintf(os.Stderr, "test %d (%s): expected a valid number but got -1\ninput:\n%s\n", idx+1, tc.desc, input)
			os.Exit(1)
		}
		if err := verifyNumber(tc.a, gotNum); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s): invalid number: %v\ninput:\n%s\nnumber:\n%s\n", idx+1, tc.desc, err, input, gotNum)
			os.Exit(1)
		}
	}

	if executed == 0 {
		fmt.Println("No usable tests were executed.")
		return
	}
	fmt.Printf("All %d tests passed.\n", executed)
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1120E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "1120E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(path, input string) (string, string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseOutput(out string) (string, bool, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return "", false, fmt.Errorf("empty output")
	}
	if tokens[0] == "-1" {
		if len(tokens) != 1 {
			return "", false, fmt.Errorf("unexpected tokens after -1")
		}
		return "", true, nil
	}
	if len(tokens) != 1 {
		return "", false, fmt.Errorf("expected single integer, got %d tokens", len(tokens))
	}
	return tokens[0], false, nil
}

func verifyNumber(a int, num string) error {
	if len(num) == 0 {
		return fmt.Errorf("empty number")
	}
	if len(num) > digitLimit {
		return fmt.Errorf("number exceeds %d digits", digitLimit)
	}
	if num[0] == '0' {
		return fmt.Errorf("leading zero not allowed")
	}

	digits := make([]int, len(num))
	sumN := 0
	for i := 0; i < len(num); i++ {
		ch := num[i]
		if ch < '0' || ch > '9' {
			return fmt.Errorf("invalid character %q", ch)
		}
		val := int(ch - '0')
		digits[len(num)-1-i] = val
		sumN += val
	}

	if sumN%a != 0 {
		return fmt.Errorf("digit sum %d not divisible by %d", sumN, a)
	}
	target := sumN / a

	sumProd := 0
	carry := 0
	for _, d := range digits {
		prod := d*a + carry
		sumProd += prod % 10
		carry = prod / 10
	}
	for carry > 0 {
		sumProd += carry % 10
		carry /= 10
	}
	if sumProd != target {
		return fmt.Errorf("S(a*n)=%d but expected %d", sumProd, target)
	}
	return nil
}

func generateTests() []testCase {
	tests := []testCase{
		{a: 2, desc: "sample-2"},
		{a: 3, desc: "sample-3"},
		{a: 10, desc: "sample-10"},
		{a: 4, desc: "small-4"},
		{a: 5, desc: "small-5"},
		{a: 6, desc: "small-6"},
		{a: 7, desc: "small-7"},
		{a: 8, desc: "small-8"},
		{a: 9, desc: "small-9"},
		{a: 12, desc: "mid-12"},
		{a: 15, desc: "mid-15"},
		{a: 16, desc: "mid-16"},
		{a: 18, desc: "mid-18"},
		{a: 20, desc: "mid-20"},
		{a: 25, desc: "power-of-5"},
		{a: 30, desc: "mid-30"},
		{a: 40, desc: "mid-40"},
		{a: 50, desc: "mid-50"},
		{a: 60, desc: "mid-60"},
		{a: 75, desc: "mid-75"},
		{a: 100, desc: "mid-100"},
		{a: 125, desc: "power-of-5^3"},
		{a: 150, desc: "mid-150"},
		{a: 200, desc: "mid-200"},
	}

	rng := rand.New(rand.NewSource(112000))
	for i := 0; i < 25; i++ {
		val := rng.Intn(199) + 2 // range [2, 200]
		tests = append(tests, testCase{
			a:    val,
			desc: fmt.Sprintf("random-%d", i+1),
		})
	}
	return tests
}
