package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	refSource     = "2000-2999/2000-2099/2080-2089/2084/2084F.go"
	totalNLimit   = 500000
	defaultTimout = 25 * time.Second
)

type testCase struct {
	n int
	a []int
	c []int
}

type parsedAnswer struct {
	isMinusOne bool
	perm       []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutputs(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutputs(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		if refAns[i].isMinusOne {
			if !candAns[i].isMinusOne {
				fmt.Fprintf(os.Stderr, "test %d: expected -1 but candidate provided an answer\ninput:\n%s", i+1, singleInput(tc))
				os.Exit(1)
			}
			continue
		}
		if candAns[i].isMinusOne {
			fmt.Fprintf(os.Stderr, "test %d: candidate output -1 but a solution exists\ninput:\n%s", i+1, singleInput(tc))
			os.Exit(1)
		}
		if err := validateSolution(tc, candAns[i].perm); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid solution: %v\ninput:\n%s", i+1, err, singleInput(tc))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2084F-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2084F")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, out.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func buildTests() []testCase {
	tests := deterministicTests()
	total := 0
	for _, tc := range tests {
		total += tc.n
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < totalNLimit {
		n := rng.Intn(60000) + 2
		if total+n > totalNLimit {
			n = totalNLimit - total
		}
		if n < 2 {
			break
		}
		a := randPermutation(rng, n)
		c := make([]int, n)
		prefilled := rng.Intn(n + 1)
		usedIdx := make([]int, n)
		for i := 0; i < n; i++ {
			usedIdx[i] = i
		}
		rng.Shuffle(len(usedIdx), func(i, j int) { usedIdx[i], usedIdx[j] = usedIdx[j], usedIdx[i] })
		for i := 0; i < prefilled; i++ {
			idx := usedIdx[i]
			if rng.Intn(4) == 0 {
				continue
			}
			c[idx] = a[idx]
		}
		tests = append(tests, testCase{n: n, a: a, c: c})
		total += n
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, a: []int{2, 1}, c: []int{0, 0}},
		{n: 3, a: []int{3, 1, 2}, c: []int{0, 0, 0}},
		{n: 4, a: []int{1, 3, 2, 4}, c: []int{1, 0, 0, 0}},
		{n: 5, a: []int{3, 1, 5, 2, 4}, c: []int{0, 3, 0, 0, 4}},
		{n: 6, a: []int{6, 5, 4, 3, 2, 1}, c: []int{0, 0, 0, 0, 0, 0}},
		{n: 7, a: []int{2, 7, 1, 4, 5, 6, 3}, c: []int{0, 0, 1, 0, 0, 0, 3}},
		{n: 8, a: []int{5, 3, 1, 8, 2, 4, 7, 6}, c: []int{0, 0, 1, 0, 0, 0, 0, 0}},
	}
}

func randPermutation(rng *rand.Rand, n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { p[i], p[j] = p[j], p[i] })
	return p
}

func buildInput(tests []testCase) string {
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
		for i, v := range tc.c {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", tc.n)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimout)
	defer cancel()
	cmd := commandFor(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return outBuf.String(), nil
}

func commandFor(ctx context.Context, path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.CommandContext(ctx, "go", "run", path)
	case ".py":
		return exec.CommandContext(ctx, "python3", path)
	default:
		return exec.CommandContext(ctx, path)
	}
}

func parseOutputs(out string, tests []testCase) ([]parsedAnswer, error) {
	fields := strings.Fields(out)
	idx := 0
	res := make([]parsedAnswer, 0, len(tests))
	for _, tc := range tests {
		if idx >= len(fields) {
			return nil, fmt.Errorf("insufficient tokens for outputs")
		}
		if fields[idx] == "-1" {
			res = append(res, parsedAnswer{isMinusOne: true})
			idx++
			continue
		}
		if len(fields)-idx < tc.n {
			return nil, fmt.Errorf("not enough numbers for a permutation of length %d", tc.n)
		}
		perm := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			v, err := strconv.Atoi(fields[idx+i])
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[idx+i])
			}
			perm[i] = v
		}
		idx += tc.n
		res = append(res, parsedAnswer{perm: perm})
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra tokens in output")
	}
	return res, nil
}

func validateSolution(tc testCase, perm []int) error {
	n := tc.n
	if len(perm) != n {
		return fmt.Errorf("expected %d elements, got %d", n, len(perm))
	}
	seen := make([]bool, n+1)
	posOfVal := make([]int, n+1)
	for i, v := range perm {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range at position %d", v, i+1)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
		posOfVal[v] = i + 1
		if tc.c[i] != 0 && tc.c[i] != v {
			return fmt.Errorf("position %d must be %d per constraints, got %d", i+1, tc.c[i], v)
		}
	}

	L, R := computeBounds(tc.a)
	for val := 1; val <= n; val++ {
		pos := posOfVal[val]
		if pos < L[val] || pos > R[val] {
			return fmt.Errorf("value %d placed at %d outside allowed range [%d,%d]", val, pos, L[val], R[val])
		}
	}
	return nil
}

func computeBounds(a []int) ([]int, []int) {
	n := len(a)
	L := make([]int, n+1)
	R := make([]int, n+1)
	fen := newFenwick(n)
	for i := 0; i < n; i++ {
		val := a[i]
		L[val] = 1 + fen.sum(val-1)
		fen.add(val, 1)
	}
	fen = newFenwick(n)
	for i := n - 1; i >= 0; i-- {
		val := a[i]
		greater := fen.sum(n) - fen.sum(val)
		R[val] = n - greater
		fen.add(val, 1)
	}
	return L, R
}

type fenwick struct {
	tree []int
	n    int
}

func newFenwick(n int) *fenwick {
	return &fenwick{tree: make([]int, n+2), n: n}
}

func (f *fenwick) add(idx, delta int) {
	for idx <= f.n {
		f.tree[idx] += delta
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	res := 0
	for idx > 0 {
		res += f.tree[idx]
		idx -= idx & -idx
	}
	return res
}
