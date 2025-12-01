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
)

const refSourceF = "./2004F.go"

type caseBundle struct {
	arrays [][]int
}

func (cb caseBundle) input() string {
	var b strings.Builder
	fmt.Fprintln(&b, len(cb.arrays))
	for _, arr := range cb.arrays {
		fmt.Fprintln(&b, len(arr))
		for i, v := range arr {
			if i > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(strconv.Itoa(v))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	var candidate string
	if len(os.Args) == 2 {
		candidate = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		candidate = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for idx, tc := range tests {
		input := tc.input()

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVals, err := parseOutputs(refOut, len(tc.arrays))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVals, err := parseOutputs(candOut, len(tc.arrays))
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d case %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, i+1, input, refVals[i], candVals[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2004F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceF))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
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

func buildTests() []caseBundle {
	var tests []caseBundle

	// Statement sample.
	tests = append(tests, caseBundle{arrays: [][]int{
		{2, 1, 3},
		{1, 1, 1, 1},
		{4, 2, 3, 1, 5},
		{1, 2, 1, 2},
	}})

	// Small edge shapes.
	tests = append(tests, caseBundle{arrays: [][]int{
		{1},
		{7},
		{1, 1},
		{2, 3},
		{3, 1, 1},
	}})

	// Palindromic, near-palindromic, and high-value mixes.
	tests = append(tests, caseBundle{arrays: [][]int{
		{5, 4, 3, 4, 5},
		{10, 1, 10, 1, 10},
		{100000, 99999, 100000},
		{8, 8, 8, 8},
	}})

	// Randomized small bundles (total n well under 2000).
	rng := rand.New(rand.NewSource(2004))
	for i := 0; i < 5; i++ {
		var arrays [][]int
		total := 0
		for len(arrays) < 6 && total < 120 {
			n := rng.Intn(10) + 1
			if total+n > 120 {
				break
			}
			arrays = append(arrays, randomArray(rng, n))
			total += n
		}
		tests = append(tests, caseBundle{arrays: arrays})
	}

	// Larger stress-style bundle with total length near the limit.
	tests = append(tests, bigBundle(rng, 2000))

	return tests
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := range arr {
		// Bias towards small numbers with occasional spikes to 1e5.
		if rng.Intn(6) == 0 {
			arr[i] = rng.Intn(100000) + 1
		} else {
			arr[i] = rng.Intn(9) + 1
		}
	}
	return arr
}

func bigBundle(rng *rand.Rand, total int) caseBundle {
	var arrays [][]int
	remaining := total

	// Mix of one large case and several medium ones to hit the sum limit.
	mainLen := remaining / 2
	arrays = append(arrays, randomArray(rng, mainLen))
	remaining -= mainLen

	for remaining > 0 {
		n := rng.Intn(300) + 1
		if n > remaining {
			n = remaining
		}
		arrays = append(arrays, randomArray(rng, n))
		remaining -= n
	}
	return caseBundle{arrays: arrays}
}
