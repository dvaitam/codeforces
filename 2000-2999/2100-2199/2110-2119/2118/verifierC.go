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

const refSource = "2000-2999/2100-2199/2110-2119/2118/2118C.go"

type testCase struct {
	n   int
	k   uint64
	arr []uint64
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := args[len(args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refVals, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candVals, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if refVals[i] != candVals[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "wrong answer on case %d (n=%d k=%d): expected %d got %d\n", i+1, tc.n, tc.k, refVals[i], candVals[i])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2118C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func buildTests() []testCase {
	var tcs []testCase
	add := func(n int, k uint64, arr []uint64) {
		tcs = append(tcs, testCase{n: n, k: k, arr: arr})
	}

	// Deterministic edge cases.
	add(5, 2, []uint64{0, 1, 7, 2, 4}) // sample-like
	add(5, 3, []uint64{0, 1, 7, 2, 4})
	add(1, 0, []uint64{3})                  // no ops
	add(1, 5, []uint64{0})                  // single element growth
	add(2, 0, []uint64{1, 2})               // already set
	add(3, 10, []uint64{0, 0, 0})           // many small increments
	add(4, 6, []uint64{15, 15, 15, 15})     // no benefit
	add(4, 1, []uint64{8, 0, 0, 0})         // choose best target
	add(6, 20, []uint64{5, 6, 7, 8, 9, 10}) // mixed bits
	add(8, 50, []uint64{1, 3, 7, 15, 31, 63, 127, 255})
	add(10, 100, []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	rng := rand.New(rand.NewSource(2118))
	for len(tcs) < 40 {
		n := rng.Intn(60) + 1
		if len(tcs)%8 == 0 {
			n = rng.Intn(400) + 200
		}
		arr := make([]uint64, n)
		for i := 0; i < n; i++ {
			if rng.Intn(5) == 0 {
				arr[i] = uint64(rng.Intn(16))
			} else {
				arr[i] = uint64(rng.Intn(1_000_000_000))
			}
		}
		k := uint64(rng.Intn(200000) + 1)
		if rng.Intn(6) == 0 {
			k = 0
		}
		add(n, k, arr)
	}

	// Stress near limits but with moderate k to keep runtime reasonable.
	largeN := 5000
	largeArr := make([]uint64, largeN)
	for i := 0; i < largeN; i++ {
		largeArr[i] = uint64(rng.Intn(1_000_000_000))
	}
	add(largeN, 100000, largeArr)

	return tcs
}

func buildInput(tcs []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatUint(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
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
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseAnswers(out string, t int) ([]uint64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	res := make([]uint64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseUint(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}
