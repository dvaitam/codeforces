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
	refSource   = "./2132E.go"
	totalNLimit = 200000
	totalMLimit = 200000
	totalQLimit = 100000
	maxVal      = 1_000_000_000
	defaultTime = 20 * time.Second
)

type query struct {
	x int
	y int
	z int
}

type testCase struct {
	n  int
	m  int
	q  int
	a  []int
	b  []int
	qs []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
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

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "output length mismatch\nreference answers: %d\ncandidate answers: %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer at global output %d: expected %d got %d\ninput fragment:\n%s", i+1, refAns[i], candAns[i], testForOutput(tests, i))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2132E-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2132E")
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
	totalN, totalM, totalQ := 0, 0, 0
	for _, tc := range tests {
		totalN += tc.n
		totalM += tc.m
		totalQ += tc.q
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for totalN < totalNLimit && totalM < totalMLimit && totalQ < totalQLimit {
		n := rng.Intn(5000) + 1
		m := rng.Intn(5000) + 1
		q := rng.Intn(2000) + 1
		if totalN+n > totalNLimit {
			n = totalNLimit - totalN
		}
		if totalM+m > totalMLimit {
			m = totalMLimit - totalM
		}
		if totalQ+q > totalQLimit {
			q = totalQLimit - totalQ
		}
		if n <= 0 || m <= 0 || q <= 0 {
			break
		}
		a := make([]int, n)
		b := make([]int, m)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(maxVal) + 1
		}
		for i := 0; i < m; i++ {
			b[i] = rng.Intn(maxVal) + 1
		}
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			x := rng.Intn(n + 1)
			y := rng.Intn(m + 1)
			maxZ := x + y
			z := 0
			if maxZ > 0 {
				z = rng.Intn(maxZ + 1)
			}
			qs[i] = query{x: x, y: y, z: z}
		}
		tests = append(tests, testCase{n: n, m: m, q: q, a: a, b: b, qs: qs})
		totalN += n
		totalM += m
		totalQ += q
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1, m: 1, q: 1,
			a:  []int{5},
			b:  []int{10},
			qs: []query{{x: 1, y: 1, z: 1}},
		},
		{
			n: 2, m: 2, q: 3,
			a:  []int{1, 2},
			b:  []int{3, 4},
			qs: []query{{x: 1, y: 1, z: 2}, {x: 2, y: 0, z: 1}, {x: 0, y: 2, z: 2}},
		},
		{
			n: 3, m: 4, q: 2,
			a:  []int{10, 20, 30},
			b:  []int{5, 5, 5, 5},
			qs: []query{{x: 2, y: 2, z: 3}, {x: 3, y: 1, z: 4}},
		},
	}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.q)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.qs {
			fmt.Fprintf(&sb, "%d %d %d\n", qu.x, qu.y, qu.z)
		}
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", tc.n, tc.m, tc.q)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, qu := range tc.qs {
		fmt.Fprintf(&sb, "%d %d %d\n", qu.x, qu.y, qu.z)
	}
	return sb.String()
}

func testForOutput(tests []testCase, idx int) string {
	for _, tc := range tests {
		if idx < tc.q {
			return singleInput(tc)
		}
		idx -= tc.q
	}
	return ""
}

func runProgram(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTime)
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

func parseOutputs(out string, tests []testCase) ([]int64, error) {
	fields := strings.Fields(out)
	totalOut := 0
	for _, tc := range tests {
		totalOut += tc.q
	}
	if len(fields) != totalOut {
		return nil, fmt.Errorf("expected %d numbers, got %d", totalOut, len(fields))
	}
	res := make([]int64, totalOut)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		res[i] = v
	}
	return res, nil
}
