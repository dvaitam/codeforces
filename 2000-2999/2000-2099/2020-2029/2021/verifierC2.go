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
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2020-2029/2021/2021C2.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2021C2.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n, m, q int
	a       []int
	b       []int
	updates [][2]int
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 4, m: 2, q: 2,
			a: []int{1, 2, 3, 4},
			b: []int{1, 1},
			updates: [][2]int{
				{1, 2},
				{1, 1},
			},
		},
		{
			n: 3, m: 2, q: 2,
			a: []int{1, 2, 3},
			b: []int{2, 3},
			updates: [][2]int{
				{1, 2},
				{2, 1},
			},
		},
	}
}

func randomPermutation(rng *rand.Rand, n int) []int {
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	return p
}

func randomArray(rng *rand.Rand, n, limit int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(limit) + 1
	}
	return arr
}

func randomUpdates(rng *rand.Rand, q, m, n int) [][2]int {
	updates := make([][2]int, q)
	for i := 0; i < q; i++ {
		updates[i] = [2]int{rng.Intn(m) + 1, rng.Intn(n) + 1}
	}
	return updates
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	totalM := 0
	totalQ := 0
	for len(tests) < 80 && totalN < 200000 && totalM < 200000 && totalQ < 200000 {
		n := rng.Intn(30) + 1
		m := rng.Intn(30) + 1
		q := rng.Intn(30)
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		if totalM+m > 200000 {
			m = 200000 - totalM
		}
		if totalQ+q > 200000 {
			q = 200000 - totalQ
		}
		if n <= 0 || m <= 0 || q < 0 {
			break
		}
		tests = append(tests, testCase{
			n:       n,
			m:       m,
			q:       q,
			a:       randomPermutation(rng, n),
			b:       randomArray(rng, m, n),
			updates: randomUpdates(rng, q, m, n),
		})
		totalN += n
		totalM += m
		totalQ += q
	}
	return tests
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.q))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, upd := range tc.updates {
			sb.WriteString(fmt.Sprintf("%d %d\n", upd[0], upd[1]))
		}
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2021C2-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2021C2")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutputs(out string, total int) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	res := make([]string, 0, total)
	for scanner.Scan() {
		res = append(res, strings.ToUpper(scanner.Text()))
	}
	if len(res) != total {
		return nil, fmt.Errorf("expected %d tokens, got %d", total, len(res))
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

	totalOutputs := 0
	for _, tc := range tests {
		totalOutputs += tc.q + 1
	}

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutputs(refOut, totalOutputs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, totalOutputs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	idx := 0
	for ti, tc := range tests {
		for state := 0; state < tc.q+1; state++ {
			if expected[idx] != got[idx] {
				fmt.Fprintf(os.Stderr, "test %d state %d mismatch: expected %s got %s\n", ti+1, state+1, expected[idx], got[idx])
				os.Exit(1)
			}
			idx++
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
