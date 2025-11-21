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

type testCase struct {
	desc  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
		expTokens, err := runAndNormalize(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		gotTokens, err := runAndNormalize(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		if len(expTokens) != len(gotTokens) {
			fmt.Fprintf(os.Stderr, "test %d (%s): token count mismatch\nexpected: %v\ngot: %v\n", i+1, tc.desc, expTokens, gotTokens)
			os.Exit(1)
		}
		for j := range expTokens {
			if expTokens[j] != gotTokens[j] {
				fmt.Fprintf(os.Stderr, "test %d (%s): mismatch at answer %d\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.desc, j+1, tc.input, expTokens[j], gotTokens[j])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2091D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	src := filepath.Clean("2000-2999/2000-2099/2090-2099/2091/2091D.go")
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runAndNormalize(path, input string) ([]string, error) {
	out, err := runProgram(path, input)
	if err != nil {
		return nil, err
	}
	return strings.Fields(out), nil
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

func buildTests() []testCase {
	tests := []testCase{
		{
			desc:  "sample",
			input: "5\n3 4 7\n5 5 5\n1 13 2\n2 4 7\n1 5 4\n",
		},
		{
			desc:  "single_row",
			input: formatSingleCase([]triple{{1, 10, 7}, {1, 10, 10}, {1, 1, 1}}),
		},
		{
			desc:  "single_col",
			input: formatSingleCase([]triple{{10, 1, 7}, {10, 1, 10}}),
		},
		{
			desc:  "k_equals_nm",
			input: formatSingleCase([]triple{{100, 100, 10000}, {2, 2, 4}, {3, 3, 9}}),
		},
		{
			desc:  "k_small",
			input: formatSingleCase([]triple{{100, 1000000000, 1}, {1000000000, 1000000000, 2}}),
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 60; i++ {
		cnt := rng.Intn(5) + 1
		cases := make([]triple, cnt)
		for j := 0; j < cnt; j++ {
			n := randRange(rng, 1, 1_000_000_000)
			m := randRange(rng, 1, 1_000_000_000)
			maxK := n * m
			k := randRange(rng, 1, maxK)
			cases[j] = triple{n, m, k}
		}
		tests = append(tests, testCase{
			desc:  fmt.Sprintf("rand-%d", i+1),
			input: formatSingleCase(cases),
		})
	}
	return tests
}

type triple struct {
	n int64
	m int64
	k int64
}

func formatSingleCase(cases []triple) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", c.n, c.m, c.k))
	}
	return sb.String()
}

func randRange(rng *rand.Rand, l, r int64) int64 {
	if l == r {
		return l
	}
	return l + rng.Int63n(r-l+1)
}
