package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	desc  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, tc := range tests {
		expect, err := runAndParse(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		got, err := runAndParse(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		if expect.Cmp(got) != 0 {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s)\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.desc, tc.input, expect.String(), got.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1662G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	src := filepath.Clean("./1662G.go")
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runAndParse(bin, input string) (*big.Int, error) {
	out, err := runProgram(bin, input)
	if err != nil {
		return nil, err
	}
	ans, err := parseBigInt(out)
	if err != nil {
		return nil, err
	}
	return ans, nil
}

func runProgram(path, input string) (string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func parseBigInt(out string) (*big.Int, error) {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	if len(tokens) > 1 {
		return nil, fmt.Errorf("expected single integer, got: %q", out)
	}
	val := new(big.Int)
	if _, ok := val.SetString(tokens[0], 10); !ok {
		return nil, fmt.Errorf("cannot parse integer %q", tokens[0])
	}
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		makeTest("n2_line", 2, []int{1}),
		makeTest("n3_line", 3, []int{1, 2}),
		makeTest("n5_sample_like", 5, []int{1, 1, 2, 2}),
		makeTest("n6_star", 6, []int{1, 1, 1, 1, 1}),
		makeTest("n6_chain", 6, []int{1, 2, 3, 4, 5}),
		makeTest("n8_balanced", 8, []int{1, 1, 2, 2, 3, 3, 4}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		n := rng.Intn(90) + 10 // 10..99
		tests = append(tests, randomTest(fmt.Sprintf("rand-small-%d", i+1), n, rng))
	}
	largeSizes := []int{200, 1000, 5000}
	for idx, n := range largeSizes {
		tests = append(tests, randomTest(fmt.Sprintf("rand-large-%d", idx+1), n, rng))
	}
	return tests
}

func makeTest(desc string, n int, parents []int) testCase {
	return testCase{desc: desc, input: formatInput(n, parents)}
}

func randomTest(desc string, n int, rng *rand.Rand) testCase {
	parents := make([]int, n-1)
	for i := 2; i <= n; i++ {
		parents[i-2] = rng.Intn(i-1) + 1
	}
	return makeTest(desc, n, parents)
}

func formatInput(n int, parents []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, p := range parents {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p))
	}
	sb.WriteByte('\n')
	return sb.String()
}
