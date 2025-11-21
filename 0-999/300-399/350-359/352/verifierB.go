package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type pair struct {
	x int
	d int
}

func solveRef(input string) ([]pair, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return nil, err
	}
	pos := make(map[int][]int)
	for i := 1; i <= n; i++ {
		var v int
		if _, err := fmt.Fscan(reader, &v); err != nil {
			return nil, err
		}
		pos[v] = append(pos[v], i)
	}
	keys := make([]int, 0, len(pos))
	for k := range pos {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var res []pair
	for _, k := range keys {
		p := pos[k]
		if len(p) == 1 {
			res = append(res, pair{x: k, d: 0})
			continue
		}
		diff := p[1] - p[0]
		ok := true
		for i := 2; i < len(p); i++ {
			if p[i]-p[i-1] != diff {
				ok = false
				break
			}
		}
		if ok {
			res = append(res, pair{x: k, d: diff})
		}
	}
	return res, nil
}

func parseOutput(out string) (int, []pair, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, nil, fmt.Errorf("failed to read count: %v", err)
	}
	pairs := make([]pair, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &pairs[i].x, &pairs[i].d); err != nil {
			return 0, nil, fmt.Errorf("failed to read pair #%d: %v", i+1, err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return 0, nil, fmt.Errorf("extraneous output detected (e.g., %q)", extra)
	}
	return t, pairs, nil
}

func validateOutput(expect []pair, candidate []pair) error {
	if len(candidate) != len(expect) {
		return fmt.Errorf("expected %d entries but got %d", len(expect), len(candidate))
	}
	seen := make(map[int]bool)
	prevX := -1
	for idx, p := range candidate {
		if seen[p.x] {
			return fmt.Errorf("duplicate x value %d", p.x)
		}
		seen[p.x] = true
		if idx > 0 && p.x <= prevX {
			return fmt.Errorf("x values must be in strictly increasing order")
		}
		prevX = p.x
	}
	expectMap := make(map[int]int, len(expect))
	for _, p := range expect {
		expectMap[p.x] = p.d
	}
	for _, p := range candidate {
		d, ok := expectMap[p.x]
		if !ok {
			return fmt.Errorf("unexpected x value %d", p.x)
		}
		if p.d != d {
			return fmt.Errorf("x=%d expected difference %d but got %d", p.x, d, p.d)
		}
	}
	if len(seen) != len(expect) {
		return fmt.Errorf("missing some expected values")
	}
	return nil
}

type testCase struct {
	name  string
	input string
}

func makeCase(name string, arr []int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return testCase{name: name, input: sb.String()}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(352))
	var tests []testCase
	gen := func(prefix string, cases, maxN, maxVal int) {
		for i := 0; i < cases; i++ {
			n := rng.Intn(maxN) + 1
			arr := make([]int, n)
			for j := 0; j < n; j++ {
				arr[j] = rng.Intn(maxVal) + 1
			}
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), arr))
		}
	}
	gen("small", 120, 8, 5)
	gen("medium", 120, 40, 50)
	gen("large", 60, 200, 100000)
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_element", []int{5}),
		makeCase("all_same", []int{2, 2, 2, 2}),
		makeCase("two_values_ap", []int{1, 2, 1, 2, 1, 2}),
		makeCase("non_ap_positions", []int{3, 3, 3, 4, 3}),
		makeCase("mixed", []int{1, 3, 1, 3, 1, 3, 5}),
		makeCase("increasing_sequence", []int{1, 2, 3, 4, 5}),
		makeCase("sparse_hits", []int{7, 1, 7, 2, 7, 3, 7}),
		makeCase("large_value_gap", []int{100000, 1, 100000, 2, 100000}),
	}
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to build reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		output, runErr := runCandidate(bin, tc.input)
		if runErr != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, runErr, tc.input)
			os.Exit(1)
		}
		t, pairs, parseErr := parseOutput(output)
		if parseErr != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, parseErr, tc.input, output)
			os.Exit(1)
		}
		if t != len(pairs) {
			fmt.Printf("test %d (%s) declared count %d but parsed %d pairs\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, t, len(pairs), tc.input, output)
			os.Exit(1)
		}
		if err := validateOutput(expect, pairs); err != nil {
			fmt.Printf("test %d (%s) failed validation: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, output)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
