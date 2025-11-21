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
	"time"
)

const refSource = "2011G.go"

type testCase struct {
	name         string
	input        string
	caseLengths  []int
	totalOutputs int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate_binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refParsed, err := parseOutput(refOut, tc.caseLengths, tc.totalOutputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candParsed, err := parseOutput(candOut, tc.caseLengths, tc.totalOutputs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for caseIdx := range refParsed {
			exp := refParsed[caseIdx]
			got := candParsed[caseIdx]
			for pos := range exp {
				if got[pos] != exp[pos] {
					fmt.Fprintf(os.Stderr, "test %d (%s) mismatch in case %d position %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
						idx+1, tc.name, caseIdx+1, pos+1, exp[pos], got[pos], tc.input, refOut, candOut)
					os.Exit(1)
				}
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2011G-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, caseLengths []int, total int) ([][]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) != total {
		return nil, fmt.Errorf("expected %d integers, got %d", total, len(tokens))
	}
	res := make([][]int64, len(caseLengths))
	pos := 0
	for i, need := range caseLengths {
		cur := make([]int64, need)
		for j := 0; j < need; j++ {
			val, err := strconv.ParseInt(tokens[pos], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q: %v", tokens[pos], err)
			}
			if val < 0 {
				return nil, fmt.Errorf("negative answer %d", val)
			}
			cur[j] = val
			pos++
		}
		res[i] = cur
	}
	return res, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newManualTest("tiny_desc", [][]int{
			{3, 2, 1},
		}),
		newManualTest("tiny_inc", [][]int{
			{1, 2, 3, 4},
		}),
		newManualTest("two_cases", [][]int{
			{4, 2, 1, 3},
			{5, 2, 4, 1, 3},
		}),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng, i))
	}
	return tests
}

func newManualTest(name string, perms [][]int) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(perms)))
	sb.WriteByte('\n')
	caseLens := make([]int, 0, len(perms))
	for _, perm := range perms {
		n := len(perm)
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for i, v := range perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		caseLens = append(caseLens, n-1)
	}
	total := 0
	for _, l := range caseLens {
		total += l
	}
	return testCase{
		name:         name,
		input:        sb.String(),
		caseLengths:  caseLens,
		totalOutputs: total,
	}
}

func randomTest(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	caseLens := make([]int, t)
	total := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(12) + 2
		if rng.Intn(5) == 0 {
			n = rng.Intn(50) + 2
		}
		perm := randPerm(rng, n)
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j, v := range perm {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		caseLens[i] = n - 1
		total += n - 1
	}
	return testCase{
		name:         fmt.Sprintf("random_%d", idx+1),
		input:        sb.String(),
		caseLengths:  caseLens,
		totalOutputs: total,
	}
}

func randPerm(rng *rand.Rand, n int) []int {
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		perm[i], perm[j] = perm[j], perm[i]
	})
	return perm
}
