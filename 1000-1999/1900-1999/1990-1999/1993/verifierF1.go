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

const refSource = "1000-1999/1900-1999/1990-1999/1993/1993F1.go"

type testCase struct {
	name  string
	input string
}

func main() {
	candPath, ok := parseBinaryArg()
	if !ok {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/binary")
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
		tmp, err := os.CreateTemp("", "verifier1993F1-*")
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
	aTokens := strings.Fields(a)
	bTokens := strings.Fields(b)
	if len(aTokens) != len(bTokens) {
		return false
	}
	for i := range aTokens {
		if aTokens[i] != bTokens[i] {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	tests := []testCase{
		sampleTest(),
		edgeTest(),
	}
	tests = append(tests, randomTest("random_small", 10, 200)...)
	tests = append(tests, randomTest("random_medium", 10, 2000)...)
	tests = append(tests, randomTest("random_large", 5, 10000)...)
	return tests
}

func sampleTest() testCase {
	return testCase{
		name: "sample",
		input: strings.TrimSpace(`5
2 2 2 2
UR
4 2 1 1
LLDD
6 3 3 1
RLRRRL
5 5 3 3
RUURD
7 5 3 4
RRDLUUU
`) + "\n",
	}
}

func edgeTest() testCase {
	return testCase{
		name: "edge",
		input: strings.TrimSpace(`3
1 1 1 1
U
1 1 1000000 1000000
R
3 3 2 2
LUD
`) + "\n",
	}
}

func randomTest(name string, batchCount, maxN int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for b := 0; b < batchCount; b++ {
		var sb strings.Builder
		caseCnt := rng.Intn(4) + 1
		fmt.Fprintf(&sb, "%d\n", caseCnt)
		totalN := 0
		for i := 0; i < caseCnt; i++ {
			n := rng.Intn(maxN) + 1
			k := rng.Intn(n) + 1
			w := rng.Intn(1000000) + 1
			h := rng.Intn(1000000) + 1
			fmt.Fprintf(&sb, "%d %d %d %d\n", n, k, w, h)
			var script strings.Builder
			for j := 0; j < n; j++ {
				switch rng.Intn(4) {
				case 0:
					script.WriteByte('L')
				case 1:
					script.WriteByte('R')
				case 2:
					script.WriteByte('U')
				default:
					script.WriteByte('D')
				}
			}
			sb.WriteString(script.String())
			sb.WriteByte('\n')
			totalN += n
			if totalN > 1_000_000 {
				break
			}
		}
		tests = append(tests, testCase{name: fmt.Sprintf("%s_%d", name, b+1), input: sb.String()})
	}
	return tests
}
