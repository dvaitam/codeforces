package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type testCase struct {
	name  string
	input string
	arr   []int64
}

const maxMoves = 1000000

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expectSolution, err := parseReferenceOutcome(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := validateCandidateOutput(candOut, tc, expectSolution); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\nsolution output:\n%s", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-341E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "341E.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func parseReferenceOutcome(output string) (bool, error) {
	reader := strings.NewReader(output)
	var token string
	if _, err := fmt.Fscan(reader, &token); err != nil {
		return false, fmt.Errorf("could not read first token: %v", err)
	}
	if token == "-1" {
		return false, nil
	}
	if _, err := strconv.Atoi(token); err != nil {
		return false, fmt.Errorf("unexpected first token %q", token)
	}
	return true, nil
}

func validateCandidateOutput(output string, tc testCase, expectSolution bool) error {
	reader := strings.NewReader(output)
	var firstToken string
	if _, err := fmt.Fscan(reader, &firstToken); err != nil {
		if expectSolution {
			return fmt.Errorf("expected a sequence of moves but could not read first token: %v", err)
		}
		return fmt.Errorf("expected -1 but output is empty or invalid: %v", err)
	}
	if firstToken == "-1" {
		if expectSolution {
			return fmt.Errorf("expected a valid sequence of moves but got -1")
		}
		if err := ensureNoExtraTokens(reader); err != nil {
			return err
		}
		return nil
	}
	if !expectSolution {
		return fmt.Errorf("expected -1 but first token is %q", firstToken)
	}
	c, err := strconv.Atoi(firstToken)
	if err != nil {
		return fmt.Errorf("first token is not a valid integer: %v", err)
	}
	if c < 0 || c > maxMoves {
		return fmt.Errorf("number of moves %d is outside [0, %d]", c, maxMoves)
	}

	arr := append([]int64(nil), tc.arr...)
	n := len(arr)
	for move := 0; move < c; move++ {
		var i, j int
		if _, err := fmt.Fscan(reader, &i, &j); err != nil {
			return fmt.Errorf("failed to read move %d: %v", move+1, err)
		}
		if i < 1 || i > n || j < 1 || j > n {
			return fmt.Errorf("move %d has indices (%d, %d) outside [1, %d]", move+1, i, j, n)
		}
		if i == j {
			return fmt.Errorf("move %d uses identical source and target box %d", move+1, i)
		}
		ai := arr[i-1]
		aj := arr[j-1]
		if ai > aj {
			return fmt.Errorf("move %d violates ai<=aj (ai=%d, aj=%d)", move+1, ai, aj)
		}
		arr[i-1] = ai + ai
		arr[j-1] = aj - ai
		if arr[j-1] < 0 {
			return fmt.Errorf("box %d became negative after move %d", j, move+1)
		}
	}
	if err := ensureNoExtraTokens(reader); err != nil {
		return err
	}

	var sumBefore, sumAfter int64
	for _, v := range tc.arr {
		sumBefore += v
	}

	positive := 0
	for idx, v := range arr {
		if v < 0 {
			return fmt.Errorf("box %d has negative candies (%d) after simulation", idx+1, v)
		}
		sumAfter += v
		if v > 0 {
			positive++
		}
	}
	if sumBefore != sumAfter {
		return fmt.Errorf("total candies changed from %d to %d", sumBefore, sumAfter)
	}
	if positive != 2 {
		return fmt.Errorf("need exactly two boxes with candies, found %d", positive)
	}
	return nil
}

func ensureNoExtraTokens(reader *strings.Reader) error {
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("unexpected extra token %q after declared moves", extra)
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(341))
	var tests []testCase

	add := func(name string, arr []int64) {
		tests = append(tests, makeTestCase(name, arr))
	}

	add("simple_three_boxes", []int64{3, 12, 6})
	add("two_boxes_ready", []int64{0, 7, 0, 5})
	add("all_zero", []int64{0, 0, 0, 0})
	add("single_positive", []int64{0, 0, 15, 0})
	add("minimal_three_large", []int64{333333, 333333, 333334})
	add("sparse_large", buildSparseArray())
	add("small_random", randomArrayWithPos(rng, 6, 100, 4))
	add("medium_random", randomArrayWithPos(rng, 50, 20000, 10))
	add("large_random", randomArrayWithPos(rng, 400, 2500, 30))
	add("max_sum_random", randomArrayWithPos(rng, 1000, 1000, 60))
	add("already_two_in_large_n", func() []int64 {
		arr := make([]int64, 20)
		arr[4] = 500
		arr[15] = 700
		return arr
	}())

	return tests
}

func makeTestCase(name string, arr []int64) testCase {
	if len(arr) < 3 {
		panic("test case must have at least 3 boxes")
	}
	var sum int64
	for _, v := range arr {
		if v < 0 {
			panic("negative value in test case")
		}
		sum += v
	}
	if sum > 1000000 {
		panic(fmt.Sprintf("sum exceeds 1e6 in test %s: %d", name, sum))
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", v)
	}
	b.WriteByte('\n')
	return testCase{
		name:  name,
		input: b.String(),
		arr:   append([]int64(nil), arr...),
	}
}

func randomArrayWithPos(rng *rand.Rand, n int, maxVal int64, minPositive int) []int64 {
	if n < 3 {
		n = 3
	}
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = rng.Int63n(maxVal + 1)
	}
	for i := 0; i < minPositive && i < n; i++ {
		if arr[i] == 0 {
			arr[i] = 1 + rng.Int63n(maxVal)
			if arr[i] == 0 {
				arr[i] = 1
			}
		}
	}
	return arr
}

func buildSparseArray() []int64 {
	arr := make([]int64, 1000)
	arr[0] = 500000
	arr[1] = 300000
	arr[2] = 150000
	arr[3] = 50000
	return arr
}
