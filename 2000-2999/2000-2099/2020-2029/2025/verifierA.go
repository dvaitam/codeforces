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

const refSource = "2025A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := generateTests()
	for idx, tc := range tests {
		expect, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		got, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, got)
			os.Exit(1)
		}

		if !equalTokens(expect, got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildBinary(path string) (string, func(), error) {
	cleanPath := filepath.Clean(path)
	if strings.HasSuffix(cleanPath, ".go") {
		tmp, err := os.CreateTemp("", "verifier2025A-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		source := filepath.Join(".", cleanPath)
		cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", nil, err
	}
	return abs, func() {}, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return stdout.String(), nil
}

func equalTokens(a, b string) bool {
	ta := strings.Fields(a)
	tb := strings.Fields(b)
	if len(ta) != len(tb) {
		return false
	}
	for i := range ta {
		if ta[i] != tb[i] {
			return false
		}
	}
	return true
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		name:  "sample",
		input: sampleInput(),
	})
	tests = append(tests, buildSingleCase("simple_equal", [][2]string{{"GARAGE", "GARAGE"}}))
	tests = append(tests, buildSingleCase("single_char", [][2]string{{"A", "B"}}))
	tests = append(tests, buildSingleCase("copy_pref", [][2]string{{"ABCDE", "AABCD"}}))
	tests = append(tests, buildSingleCase("no_copy", [][2]string{{"TRAINING", "DRAINING"}}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_small_%d", i+1), rng, 5, 10))
	}
	for i := 0; i < 4; i++ {
		tests = append(tests, randomBatch(fmt.Sprintf("random_mid_%d", i+1), rng, 8, 50))
	}
	tests = append(tests, randomBatch("random_large", rng, 10, 100))

	return tests
}

func sampleInput() string {
	return `3
GARAGE
GARAGE
FORSALE
ABCDE
AABCD
TRAINING
DRAINING
`
}

func buildSingleCase(name string, pairs [][2]string) testCase {
	var b strings.Builder
	fmt.Fprintln(&b, len(pairs))
	for _, pr := range pairs {
		fmt.Fprintln(&b, pr[0])
		fmt.Fprintln(&b, pr[1])
	}
	return testCase{name: name, input: b.String()}
}

func randomBatch(name string, rng *rand.Rand, cases int, maxLen int) testCase {
	if cases < 1 {
		cases = 1
	}
	var b strings.Builder
	fmt.Fprintln(&b, cases)
	for i := 0; i < cases; i++ {
		s := randomString(rng, rng.Intn(maxLen)+1)
		t := randomString(rng, rng.Intn(maxLen)+1)
		fmt.Fprintln(&b, s)
		fmt.Fprintln(&b, t)
	}
	return testCase{name: name, input: b.String()}
}

func randomString(rng *rand.Rand, length int) string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(letters[rng.Intn(len(letters))])
	}
	return sb.String()
}
