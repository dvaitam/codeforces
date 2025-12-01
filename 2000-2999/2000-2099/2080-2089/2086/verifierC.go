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

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2086C.go"

type testCase struct {
	input string
}

type instance struct {
	n     int
	perm  []int
	order []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		expected, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if err := compareOutputs(expected, got); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, err, tc.input, strings.TrimSpace(expected), strings.TrimSpace(got))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2086C-ref-*")
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

func compareOutputs(expected, got string) error {
	expNums, err := parseInts(expected)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}
	gotNums, err := parseInts(got)
	if err != nil {
		return fmt.Errorf("output is not a list of integers: %v", err)
	}
	if len(expNums) != len(gotNums) {
		return fmt.Errorf("expected %d integers, got %d", len(expNums), len(gotNums))
	}
	for idx := range expNums {
		if expNums[idx] != gotNums[idx] {
			return fmt.Errorf("mismatch at position %d: expected %d got %d", idx+1, expNums[idx], gotNums[idx])
		}
	}
	return nil
}

func parseInts(s string) ([]int64, error) {
	fields := strings.Fields(s)
	nums := make([]int64, len(fields))
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", f)
		}
		nums[i] = val
	}
	return nums, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20862086))
	var tests []testCase

	tests = append(tests, sampleTest())
	tests = append(tests, exhaustiveSmall())
	tests = append(tests, singleCycleCase())
	tests = append(tests, layeredCyclesCase())

	for i := 0; i < 3; i++ {
		tests = append(tests, randomBatch(rng, 60, 400, 200000))
	}

	tests = append(tests, randomBatch(rng, 80, 2000, 200000))
	tests = append(tests, maxStressCase(rng))

	return tests
}

func sampleTest() testCase {
	cases := []instance{
		{
			n:     3,
			perm:  []int{1, 2, 3},
			order: []int{3, 2, 1},
		},
		{
			n:     5,
			perm:  []int{4, 5, 3, 1, 2},
			order: []int{4, 5, 1, 3, 2},
		},
		{
			n:     7,
			perm:  []int{4, 3, 1, 2, 7, 5, 6},
			order: []int{1, 2, 3, 4, 5, 6, 7},
		},
	}
	return buildInput(cases)
}

func exhaustiveSmall() testCase {
	var cases []instance
	for n := 1; n <= 4; n++ {
		perms := permutations(n)
		for _, p := range perms {
			for _, d := range perms {
				cp := append([]int(nil), p...)
				cd := append([]int(nil), d...)
				cases = append(cases, instance{n: n, perm: cp, order: cd})
			}
		}
	}
	return buildInput(cases)
}

func singleCycleCase() testCase {
	n := 5000
	perm := make([]int, n)
	for i := 0; i < n-1; i++ {
		perm[i] = i + 2
	}
	perm[n-1] = 1
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = n - i
	}
	return buildInput([]instance{{n: n, perm: perm, order: order}})
}

func layeredCyclesCase() testCase {
	lengths := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 6}
	total := 0
	for _, v := range lengths {
		total += v
	}
	perm := make([]int, total)
	var order []int
	start := 1
	for _, ln := range lengths {
		cycle := make([]int, ln)
		for i := 0; i < ln; i++ {
			cycle[i] = start + i
		}
		for i := 0; i < ln-1; i++ {
			perm[cycle[i]-1] = cycle[i+1]
		}
		perm[cycle[ln-1]-1] = cycle[0]
		order = append(order, cycle[ln-1])
		order = append(order, cycle[:ln-1]...)
		start += ln
	}
	return buildInput([]instance{{n: total, perm: perm, order: order}})
}

func randomBatch(rng *rand.Rand, maxCases, maxN, limit int) testCase {
	var cases []instance
	sumN := 0
	targetCases := rng.Intn(maxCases) + 1
	for len(cases) < targetCases && sumN < limit {
		n := rng.Intn(maxN) + 1
		if sumN+n > limit {
			n = limit - sumN
		}
		if n == 0 {
			break
		}
		perm := randPerm1Based(n, rng)
		order := randPerm1Based(n, rng)
		cases = append(cases, instance{n: n, perm: perm, order: order})
		sumN += n
	}
	if len(cases) == 0 {
		cases = append(cases, instance{
			n:     1,
			perm:  []int{1},
			order: []int{1},
		})
	}
	return buildInput(cases)
}

func maxStressCase(rng *rand.Rand) testCase {
	const n = 200000
	perm := randPerm1Based(n, rng)
	order := randPerm1Based(n, rng)
	return buildInput([]instance{{n: n, perm: perm, order: order}})
}

func randPerm1Based(n int, rng *rand.Rand) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func permutations(n int) [][]int {
	base := make([]int, n)
	for i := 0; i < n; i++ {
		base[i] = i + 1
	}
	var res [][]int
	var gen func(int)
	gen = func(pos int) {
		if pos == n {
			cp := append([]int(nil), base...)
			res = append(res, cp)
			return
		}
		for i := pos; i < n; i++ {
			base[pos], base[i] = base[i], base[pos]
			gen(pos + 1)
			base[pos], base[i] = base[i], base[pos]
		}
	}
	gen(0)
	return res
}

func buildInput(cases []instance) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, cs := range cases {
		fmt.Fprintf(&b, "%d\n", cs.n)
		writeArray(&b, cs.perm)
		writeArray(&b, cs.order)
	}
	return testCase{input: b.String()}
}

func writeArray(b *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteByte('\n')
}
