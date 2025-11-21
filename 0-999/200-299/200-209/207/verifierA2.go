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

const (
	refSourceA2 = "0-999/200-299/200-209/207/207A2.go"
	totalLimit  = 500000
)

type testCase struct {
	n    int
	k    []int
	seqs [][]int64
}

type result struct {
	badPairs int64
	order    []pair
}

type pair struct {
	val int64
	id  int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refRes, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candRes, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if candRes[i].badPairs != refRes[i].badPairs {
			fmt.Fprintf(os.Stderr, "test %d: expected %d bad pairs got %d\n", i+1, refRes[i].badPairs, candRes[i].badPairs)
			os.Exit(1)
		}
		if len(tc.flatten()) <= totalLimit {
			if err := validateOrder(tc, candRes[i]); err != nil {
				fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
				os.Exit(1)
			}
			if countBadPairs(candRes[i].order) != candRes[i].badPairs {
				fmt.Fprintf(os.Stderr, "test %d: order bad pairs mismatch\n", i+1)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "207A2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceA2))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
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

func parseOutput(output string, tests []testCase) ([]result, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	res := make([]result, len(tests))
	idx := 0
	for i, tc := range tests {
		if idx >= len(lines) {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		line := strings.TrimSpace(lines[idx])
		idx++
		if line == "" {
			i--
			continue
		}
		val, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("test %d: invalid bad pair count", i+1)
		}
		res[i].badPairs = val
		total := tc.totalTasks()
		if total <= totalLimit {
			res[i].order = make([]pair, 0, total)
			for j := 0; j < total; j++ {
				if idx >= len(lines) {
					return nil, fmt.Errorf("test %d: missing order lines", i+1)
				}
				line = strings.TrimSpace(lines[idx])
				idx++
				if line == "" {
					j--
					continue
				}
				fields := strings.Fields(line)
				if len(fields) != 2 {
					return nil, fmt.Errorf("test %d: invalid order line", i+1)
				}
				valV, err1 := strconv.ParseInt(fields[0], 10, 64)
				idV, err2 := strconv.Atoi(fields[1])
				if err1 != nil || err2 != nil {
					return nil, fmt.Errorf("test %d: invalid values in order line", i+1)
				}
				res[i].order = append(res[i].order, pair{valV, idV})
			}
		}
	}
	return res, nil
}

func (tc testCase) totalTasks() int {
	sum := 0
	for _, k := range tc.k {
		sum += k
	}
	return sum
}

func (tc testCase) flatten() []pair {
	total := tc.totalTasks()
	ans := make([]pair, 0, total)
	for i, seq := range tc.seqs {
		for _, val := range seq {
			ans = append(ans, pair{val, i + 1})
		}
	}
	return ans
}

func validateOrder(tc testCase, res result) error {
	all := tc.flatten()
	if len(res.order) != len(all) {
		return fmt.Errorf("order length mismatch")
	}
	pos := make(map[int]int)
	for i := range tc.seqs {
		pos[i+1] = 0
	}
	for _, it := range res.order {
		if it.id < 1 || it.id > len(tc.seqs) {
			return fmt.Errorf("invalid scientist id %d", it.id)
		}
		idx := pos[it.id]
		if idx >= len(tc.seqs[it.id-1]) {
			return fmt.Errorf("sequence for scientist %d exhausted", it.id)
		}
		if tc.seqs[it.id-1][idx] != it.val {
			return fmt.Errorf("value mismatch for scientist %d", it.id)
		}
		pos[it.id]++
	}
	for id, idx := range pos {
		if idx != len(tc.seqs[id-1]) {
			return fmt.Errorf("not all tasks for scientist %d used", id)
		}
	}
	return nil
}

func countBadPairs(order []pair) int64 {
	var prev int64 = math.MinInt64
	var bad int64
	for _, item := range order {
		if prev != math.MinInt64 && item.val < prev {
			bad++
		}
		prev = item.val
	}
	return bad
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(207))
	var tests []testCase
	total := 0
	add := func(tc testCase) {
		if total+tc.totalTasks() > totalLimit {
			return
		}
		tests = append(tests, tc)
		total += tc.totalTasks()
	}

	add(buildSimpleCase(2, []int{2, 2}, [][]int64{{1, 2}, {3, 4}}))
	add(buildRandomCase(rng, 2, []int{3, 3}, []int64{10, 2, 3, 1000}, []int64{100, 1, 999, 1000}))

	for total < totalLimit {
		n := rng.Intn(5) + 1
		k := make([]int, n)
		seqs := make([][]int64, n)
		for i := 0; i < n; i++ {
			k[i] = rng.Intn(50) + 1
			seq := make([]int64, k[i])
			val := int64(rng.Intn(1000))
			for j := 0; j < k[i]; j++ {
				seq[j] = val + int64(j)
			}
			seqs[i] = seq
		}
		add(testCase{n: n, k: k, seqs: seqs})
	}
	return tests
}

func buildSimpleCase(n int, k []int, seqs [][]int64) testCase {
	return testCase{n: n, k: k, seqs: seqs}
}

func buildRandomCase(rng *rand.Rand, n int, k []int, params ...[]int64) testCase {
	seqs := make([][]int64, n)
	for i := 0; i < n; i++ {
		seq := make([]int64, k[i])
		val := int64(rng.Intn(1000))
		for j := 0; j < k[i]; j++ {
			seq[j] = (val + int64(j)*int64(rng.Intn(10)+1)) % 1000
		}
		seqs[i] = seq
	}
	return testCase{n: n, k: k, seqs: seqs}
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		for i := 0; i < tc.n; i++ {
			fmt.Fprintf(&b, "%d ", tc.k[i])
			if tc.k[i] > 0 {
				fmt.Fprintf(&b, "%d 1 1 1000000000\n", tc.seqs[i][0])
			} else {
				fmt.Fprintln(&b)
			}
		}
	}
	return b.String()
}
