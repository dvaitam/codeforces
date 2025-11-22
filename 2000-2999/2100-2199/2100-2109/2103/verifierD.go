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
	refSource2103D = "2103D.go"
	refBinary2103D = "ref2103D.bin"
	maxTests       = 120
	maxTotalN      = 150000
)

type testCase struct {
	n int
	a []int
}

type node struct {
	pos int
	val int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	// Sanity-check reference solution with the validator.
	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	refPerms, err := parsePermutations(refOut, tests)
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}
	if err := validateAll(refPerms, tests); err != nil {
		fmt.Printf("reference failed validation: %v\n", err)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	candPerms, err := parsePermutations(candOut, tests)
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}
	if err := validateAll(candPerms, tests); err != nil {
		fmt.Printf("candidate failed validation: %v\n", err)
		fmt.Println("Input used:")
		fmt.Println(string(input))
		return
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2103D, refSource2103D)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2103D), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parsePermutations(out string, tests []testCase) ([][]int, error) {
	fields := strings.Fields(out)
	expect := 0
	for _, tc := range tests {
		expect += tc.n
	}
	if len(fields) != expect {
		return nil, fmt.Errorf("expected %d numbers, got %d", expect, len(fields))
	}
	perms := make([][]int, len(tests))
	idx := 0
	for t, tc := range tests {
		perm := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			val, err := strconv.Atoi(fields[idx])
			if err != nil {
				return nil, fmt.Errorf("token %d: %v", idx+1, err)
			}
			perm[i] = val
			idx++
		}
		perms[t] = perm
	}
	return perms, nil
}

func validateAll(perms [][]int, tests []testCase) error {
	if len(perms) != len(tests) {
		return fmt.Errorf("mismatched test count")
	}
	for i := range tests {
		if err := validatePermutation(perms[i], tests[i]); err != nil {
			return fmt.Errorf("test %d: %v", i+1, err)
		}
	}
	return nil
}

func validatePermutation(p []int, tc testCase) error {
	if len(p) != tc.n {
		return fmt.Errorf("expected %d values, got %d", tc.n, len(p))
	}
	seen := make([]bool, tc.n+1)
	for idx, v := range p {
		if v < 1 || v > tc.n {
			return fmt.Errorf("position %d has out of range value %d", idx+1, v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
	}

	rem, err := removalTimes(p)
	if err != nil {
		return err
	}
	for i, want := range tc.a {
		if rem[i] != want {
			return fmt.Errorf("index %d removal mismatch: expected %d, got %d", i+1, want, rem[i])
		}
	}
	return nil
}

func removalTimes(p []int) ([]int, error) {
	n := len(p)
	if n == 0 {
		return nil, fmt.Errorf("empty permutation")
	}
	res := make([]int, n)
	for i := range res {
		res[i] = -1
	}

	cur := make([]node, n)
	for i, v := range p {
		cur[i] = node{pos: i, val: v}
	}

	iter := 1
	for len(cur) > 1 {
		keep := make([]node, 0, len(cur))
		for i, nd := range cur {
			var extreme bool
			if iter%2 == 1 {
				// keep local minima
				if i == 0 {
					extreme = nd.val < cur[i+1].val
				} else if i == len(cur)-1 {
					extreme = nd.val < cur[i-1].val
				} else {
					extreme = nd.val < cur[i-1].val && nd.val < cur[i+1].val
				}
			} else {
				// keep local maxima
				if i == 0 {
					extreme = nd.val > cur[i+1].val
				} else if i == len(cur)-1 {
					extreme = nd.val > cur[i-1].val
				} else {
					extreme = nd.val > cur[i-1].val && nd.val > cur[i+1].val
				}
			}

			if extreme {
				keep = append(keep, nd)
			} else {
				res[nd.pos] = iter
			}
		}
		if len(keep) == 0 {
			return nil, fmt.Errorf("iteration %d removed all elements", iter)
		}
		cur = keep
		iter++
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2103))
	var tests []testCase
	totalN := 0

	addFromPerm := func(p []int) {
		if len(p) == 0 {
			return
		}
		if totalN+len(p) > maxTotalN || len(tests) >= maxTests {
			return
		}
		a, err := removalTimes(p)
		if err != nil {
			return
		}
		tests = append(tests, testCase{n: len(p), a: a})
		totalN += len(p)
	}

	// Deterministic cases inspired by the statement.
	addFromPerm([]int{3, 2, 1})
	addFromPerm([]int{4, 3, 5, 1, 2})
	addFromPerm([]int{6, 7, 2, 4, 3, 8, 5, 1})
	addFromPerm([]int{6, 5, 2, 1, 3, 4, 7})
	addFromPerm([]int{5, 4, 3, 2, 1})
	addFromPerm([]int{1, 2, 3, 4, 5})
	addFromPerm([]int{4, 5, 2, 3, 1})

	for len(tests) < maxTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		var n int
		switch rnd.Intn(5) {
		case 0:
			n = rnd.Intn(5) + 2
		case 1:
			n = rnd.Intn(40) + 10
		case 2:
			n = rnd.Intn(200) + 50
		case 3:
			n = rnd.Intn(2000) + 200
		default:
			n = rnd.Intn(70000) + 3000
		}
		if n > remain {
			n = remain
		}
		if n < 2 {
			n = 2
		}
		p := randomPermutation(n, rnd)
		addFromPerm(p)
	}
	return tests
}

func randomPermutation(n int, rnd *rand.Rand) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	for i := n - 1; i > 0; i-- {
		j := rnd.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}
	return p
}
