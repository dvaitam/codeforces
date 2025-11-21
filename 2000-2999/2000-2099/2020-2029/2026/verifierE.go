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

const refSource = "2000-2999/2000-2099/2020-2029/2026/2026E.go"

type testInput struct {
	name  string
	input string
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
		if !compareOutputs(refOut, candOut) {
			fmt.Fprintf(os.Stderr, "output mismatch on test %d (%s)\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, tc.name, tc.input, refOut, candOut)
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
		tmp, err := os.CreateTemp("", "verifier2026E-*")
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

func runBinary(binPath, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func compareOutputs(refOut, candOut string) bool {
	refTokens := strings.Fields(refOut)
	candTokens := strings.Fields(candOut)
	if len(refTokens) != len(candTokens) {
		return false
	}
	for i := range refTokens {
		if refTokens[i] != candTokens[i] {
			return false
		}
	}
	return true
}

func buildTests() []testInput {
	var tests []testInput
	tests = append(tests,
		makeInput("all_zero", [][]uint64{{0, 0, 0}}),
		makeInput("single_large", [][]uint64{{1<<59 + 5}}),
		makeInput("mixed_small", [][]uint64{{1, 0, 1, 2}, {7, 1, 4, 8, 14, 13, 8, 7, 6}}),
	)
	tests = append(tests, randomTests(20)...)
	return tests
}

func makeInput(name string, cases [][]uint64) testInput {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, arr := range cases {
		fmt.Fprintf(&sb, "%d\n", len(arr))
		for i, val := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
	}
	return testInput{name: name, input: sb.String()}
}

func randomTests(count int) []testInput {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testInput, 0, count)
	for i := 0; i < count; i++ {
		t := rng.Intn(5) + 1
		cases := make([][]uint64, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(100) + 1
			arr := make([]uint64, n)
			for k := 0; k < n; k++ {
				shift := rng.Intn(60)
				if shift == 59 && rng.Intn(2) == 0 {
					arr[k] = (1 << 59) | uint64(rng.Intn(1<<30))
				} else {
					arr[k] = uint64(rng.Int63()) & ((1 << 60) - 1)
				}
			}
			cases[j] = arr
		}
		tests = append(tests, makeInput(fmt.Sprintf("random_%d", i+1), cases))
	}
	return tests
}
