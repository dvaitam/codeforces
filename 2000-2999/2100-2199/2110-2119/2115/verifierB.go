package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type query struct {
	x, y, z int // 1-based
}

type testCase struct {
	n       int
	b       []int64
	queries []query
}

func simulate(a []int64, qs []query) []int64 {
	c := make([]int64, len(a))
	copy(c, a)
	for _, q := range qs {
		x, y, z := q.x-1, q.y-1, q.z-1
		v := c[x]
		if c[y] < v {
			v = c[y]
		}
		c[z] = v
	}
	return c
}

// solveRef returns a valid initial array for tc, or nil if impossible.
func solveRef(tc testCase) []int64 {
	d := make([]int64, tc.n)
	copy(d, tc.b)
	for i := len(tc.queries) - 1; i >= 0; i-- {
		o := tc.queries[i]
		req := d[o.z-1]
		d[o.z-1] = 0
		if req > d[o.x-1] {
			d[o.x-1] = req
		}
		if req > d[o.y-1] {
			d[o.y-1] = req
		}
	}
	sim := simulate(d, tc.queries)
	for i := 0; i < tc.n; i++ {
		if sim[i] != tc.b[i] {
			return nil
		}
	}
	return d
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.queries)))
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", qu.x, qu.y, qu.z))
		}
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

// parseOutputs parses the candidate's output for all tests.
// nil slice means the candidate said "-1" (impossible) for that test.
func parseOutputs(out string, tests []testCase) ([][]int64, error) {
	tokens := strings.Fields(out)
	idx := 0
	res := make([][]int64, len(tests))
	for i, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("test %d: missing output", i+1)
		}
		if tokens[idx] == "-1" {
			res[i] = nil
			idx++
		} else {
			if idx+tc.n > len(tokens) {
				return nil, fmt.Errorf("test %d: not enough numbers (got first token %q, need %d numbers)", i+1, tokens[idx], tc.n)
			}
			arr := make([]int64, tc.n)
			for j := 0; j < tc.n; j++ {
				v, err := strconv.ParseInt(tokens[idx+j], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("test %d: invalid integer %q at position %d", i+1, tokens[idx+j], j+1)
				}
				if v < 0 || v > 1_000_000_000 {
					return nil, fmt.Errorf("test %d: value %d out of range [0, 10^9]", i+1, v)
				}
				arr[j] = v
			}
			idx += tc.n
			res[i] = arr
		}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra tokens after parsing all tests")
	}
	return res, nil
}

func validate(tc testCase, ans []int64) error {
	if ans == nil {
		// Candidate says impossible; verify with solveRef.
		if solveRef(tc) != nil {
			return fmt.Errorf("candidate said -1 but a valid answer exists")
		}
		return nil
	}
	// Candidate says ans; simulate forward and check equals b.
	sim := simulate(ans, tc.queries)
	for i := 0; i < tc.n; i++ {
		if sim[i] != tc.b[i] {
			return fmt.Errorf("simulation mismatch at position %d: expected %d, got %d", i+1, tc.b[i], sim[i])
		}
	}
	return nil
}

func generateTests(rng *rand.Rand) []testCase {
	tests := []testCase{
		// Known impossible: min(c[2],c[1])→c[2], b=[1,2] requires min(...)=2>b[1]=1, impossible.
		{n: 2, b: []int64{1, 2}, queries: []query{{2, 1, 2}}},
		// n=1, no ops.
		{n: 1, b: []int64{1000000000}, queries: nil},
		// n=1, self-assign.
		{n: 1, b: []int64{42}, queries: []query{{1, 1, 1}}},
	}

	// Valid tests: start from known a, simulate to get b.
	for i := 0; i < 60; i++ {
		n := rng.Intn(10) + 1
		q := rng.Intn(15)
		a := make([]int64, n)
		for j := range a {
			a[j] = rng.Int63n(1_000_000_000) + 1
		}
		qs := make([]query, q)
		for j := range qs {
			qs[j] = query{
				x: rng.Intn(n) + 1,
				y: rng.Intn(n) + 1,
				z: rng.Intn(n) + 1,
			}
		}
		b := simulate(a, qs)
		tests = append(tests, testCase{n: n, b: b, queries: qs})
	}

	// Random tests (mix of valid and invalid).
	for i := 0; i < 40; i++ {
		n := rng.Intn(8) + 1
		q := rng.Intn(10)
		b := make([]int64, n)
		for j := range b {
			b[j] = rng.Int63n(20) + 1
		}
		qs := make([]query, q)
		for j := range qs {
			qs[j] = query{
				x: rng.Intn(n) + 1,
				y: rng.Intn(n) + 1,
				z: rng.Intn(n) + 1,
			}
		}
		tests = append(tests, testCase{n: n, b: b, queries: qs})
	}

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if strings.HasSuffix(bin, ".go") {
		tmp, err := os.CreateTemp("", "verifierB-bin-*")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create temp file: %v\n", err)
			os.Exit(1)
		}
		tmp.Close()
		defer os.Remove(tmp.Name())
		out, err := exec.Command("go", "build", "-o", tmp.Name(), bin).CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "compile error: %v\n%s", err, out)
			os.Exit(1)
		}
		bin = tmp.Name()
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := generateTests(rng)

	input := buildInput(tests)
	out, err := runBinary(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\n", err)
		os.Exit(1)
	}

	answers, err := parseOutputs(out, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "output parse error: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := validate(tc, answers[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "input:\n%s", buildInput([]testCase{tc}))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
