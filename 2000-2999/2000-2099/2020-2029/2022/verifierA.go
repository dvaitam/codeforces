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

const refSource = "2000-2999/2000-2099/2020-2029/2022/2022A.go"

type testCase struct {
	name  string
	input string
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	refBin, cleanupRef, err := buildBinary(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanupRef()

	candBin, cleanupCand, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare candidate binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanupCand()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		candOut, err := runBinary(candBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if !equalOutputs(refOut, candOut) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func parseBinaryArg() (string, bool) {
	if len(os.Args) == 2 {
		return os.Args[1], true
	}
	if len(os.Args) == 3 && os.Args[1] == "--" {
		return os.Args[2], true
	}
	return "", false
}

func buildBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "verifier2022A-*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(path))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("%v\n%s", err, out.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	abs, err := filepath.Abs(path)
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
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func equalOutputs(a, b string) bool {
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests,
		testCase{name: "sample", input: sampleInput()},
		testCase{name: "single_family", input: singleFamilyCases()},
		testCase{name: "all_pairs", input: allPairsCase()},
	)
	tests = append(tests, randomCases(20)...)
	return tests
}

func sampleInput() string {
	return strings.TrimSpace(`4
3 3
2 3 1
3 3
2 2 2
4 5
1 1 2 2
4 5
3 1 1 3
`) + "\n"
}

func singleFamilyCases() string {
	var sb strings.Builder
	sb.WriteString("3\n")
	sb.WriteString("1 1\n1\n")
	sb.WriteString("1 5\n2\n")
	sb.WriteString("1 5\n9\n")
	return sb.String()
}

func allPairsCase() string {
	var sb strings.Builder
	sb.WriteString("2\n")
	sb.WriteString("3 5\n2 2 2\n")
	sb.WriteString("5 5\n10 10 10 10 10\n")
	return sb.String()
}

func randomCases(count int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		t := rng.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for j := 0; j < t; j++ {
			n := rng.Intn(5) + 1
			r := rng.Intn(10) + n
			sb.WriteString(fmt.Sprintf("%d %d\n", n, r))
			for k := 0; k < n; k++ {
				val := rng.Intn(10) + 1
				sb.WriteString(fmt.Sprintf("%d", val))
				if k+1 < n {
					sb.WriteByte(' ')
				}
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{name: fmt.Sprintf("random_%d", i+1), input: sb.String()})
	}
	return tests
}
