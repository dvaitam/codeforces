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

const refSource = "2000-2999/2000-2099/2090-2099/2093/2093G.go"

type testCase struct {
	n   int
	k   int
	arr []int
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run verifierG.go /path/to/candidate")
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
	expected, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	if err := compare(expected, got); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2093G-ref-*")
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
	add := func(n, k int, arr []int) {
		tcs = append(tcs, testCase{n: n, k: k, arr: arr})
	}

	add(1, 0, []int{5})                           // k=0 trivial answer 1
	add(1, 1, []int{0})                           // impossible single element
	add(2, 1, []int{0, 1})                        // need both elements
	add(3, 2, []int{0, 3, 1})                     // small reachable
	add(4, 7, []int{8, 7, 0, 15})                 // xor against distant element
	add(5, 5, []int{1, 2, 3, 4, 5})               // sample-like
	add(6, 7, []int{3, 5, 1, 4, 2, 6})            // mixed small
	add(8, 10, []int{1, 1023, 1, 0, 5, 7, 9, 2})  // large xor early
	add(9, 512, []int{0, 0, 0, 0, 0, 0, 0, 0, 0}) // unreachable large k

	// Constructed alternating patterns to hit different xor values.
	add(12, 15, []int{15, 0, 15, 0, 15, 0, 1, 2, 3, 4, 5, 6})

	// Randomized tests; keep total n within constraints.
	rng := rand.New(rand.NewSource(2093))
	for len(tcs) < 25 {
		n := rng.Intn(50) + 5
		if len(tcs)%5 == 0 {
			n = rng.Intn(500) + 200
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if rng.Intn(7) == 0 {
				arr[i] = rng.Intn(16) // dense small values
			} else {
				arr[i] = rng.Intn(1_000_000_000)
			}
		}
		k := rng.Intn(1_000_000_000)
		if rng.Intn(6) == 0 {
			k = 0
		}
		add(n, k, arr)
	}

	// Large stress test near upper limits but still within total constraint.
	largeN := 195000
	largeArr := make([]int, largeN)
	for i := 0; i < largeN; i++ {
		largeArr[i] = rng.Intn(1_000_000_000)
	}
	add(largeN, 1<<29, largeArr)

	return tcs
}

func buildInput(tcs []testCase) string {
	var sb strings.Builder
	sb.Grow(32 + len(tcs)*32)
	sb.WriteString(strconv.Itoa(len(tcs)))
	sb.WriteByte('\n')
	for _, tc := range tcs {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
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

func parseAnswers(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func compare(exp, got []int64) error {
	if len(exp) != len(got) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range exp {
		if exp[i] != got[i] {
			return fmt.Errorf("test %d: expected %d, got %d", i+1, exp[i], got[i])
		}
	}
	return nil
}
