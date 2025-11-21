package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type query struct {
	typ int
	x   int
	y   int
}

type testCase struct {
	input   string
	queries []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, tc.queries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswers(gotOut, tc.queries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if len(refAns) != len(gotAns) {
			fmt.Fprintf(os.Stderr, "test %d: expected %d answers got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
				idx+1, len(refAns), len(gotAns), tc.input, refOut, gotOut)
			os.Exit(1)
		}
		for i := range refAns {
			if refAns[i] != gotAns[i] {
				fmt.Fprintf(os.Stderr, "test %d query %d mismatch: expected %d got %d\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n",
					idx+1, i+1, refAns[i], gotAns[i], tc.input, refOut, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "331B2_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "331B2.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string, queries []query) ([]int, error) {
	var answers []int
	scanner := bufio.NewScanner(strings.NewReader(out))
	for _, q := range queries {
		if q.typ != 1 {
			continue
		}
		if !scanner.Scan() {
			return nil, fmt.Errorf("expected more answers")
		}
		valStr := strings.TrimSpace(scanner.Text())
		val, err := strconv.Atoi(valStr)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", valStr)
		}
		answers = append(answers, val)
	}
	if scanner.Scan() {
		return nil, fmt.Errorf("extra output detected")
	}
	return answers, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 60)...)
	tests = append(tests, stressTests(rng)...)
	return tests
}

func manualTests() []testCase {
	n := 5
	a := []int{1, 3, 4, 2, 5}
	queries := []query{
		{typ: 1, x: 1, y: 5},
		{typ: 1, x: 3, y: 4},
		{typ: 2, x: 2, y: 3},
		{typ: 1, x: 1, y: 5},
		{typ: 2, x: 1, y: 5},
		{typ: 1, x: 1, y: 5},
	}
	return []testCase{makeTestCase(n, a, queries)}
}

func randomTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	for b := 0; b < batches; b++ {
		n := rng.Intn(50) + 2
		arr := randPermutation(rng, n)
		q := rng.Intn(100) + 1
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			if rng.Intn(2) == 0 {
				x := rng.Intn(n-1) + 1
				y := rng.Intn(n-x) + x + 1
				queries[i] = query{typ: 1, x: x, y: y}
			} else {
				x := rng.Intn(n-1) + 1
				y := rng.Intn(n-x) + x + 1
				queries[i] = query{typ: 2, x: x, y: y}
			}
		}
		tests = append(tests, makeTestCase(n, arr, queries))
	}
	return tests
}

func stressTests(rng *rand.Rand) []testCase {
	var tests []testCase
	n := 500
	arr := randPermutation(rng, n)
	queries := make([]query, 1000)
	for i := 0; i < len(queries); i++ {
		if i%3 == 0 {
			x := rng.Intn(n-1) + 1
			y := rng.Intn(n-x) + x + 1
			queries[i] = query{typ: 1, x: x, y: y}
		} else {
			x := rng.Intn(n-1) + 1
			y := rng.Intn(n-x) + x + 1
			queries[i] = query{typ: 2, x: x, y: y}
		}
	}
	tests = append(tests, makeTestCase(n, arr, queries))
	return tests
}

func makeTestCase(n int, arr []int, queries []query) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
	for _, q := range queries {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", q.typ, q.x, q.y))
	}
	return testCase{
		input:   sb.String(),
		queries: queries,
	}
}

func randPermutation(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		arr[i], arr[j] = arr[j], arr[i]
	})
	return arr
}
