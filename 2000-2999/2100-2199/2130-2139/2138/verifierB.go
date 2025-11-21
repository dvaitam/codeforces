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
	refSource   = "2000-2999/2100-2199/2130-2139/2138/2138B.go"
	totalNLimit = 500000
	totalQLimit = 500000
	defaultTime = 15 * time.Second
)

type query struct {
	l int
	r int
}

type testCase struct {
	n  int
	q  int
	a  []int
	qs []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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
		fmt.Fprintf(os.Stderr, "output length mismatch: reference %d vs candidate %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}

	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer at query %d: expected %s got %s\ninput snippet:\n%s", i+1, refAns[i], candAns[i], testForOutput(tests, i))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-2138B-*")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref2138B")
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
	totalN, totalQ := 0, 0
	for _, tc := range tests {
		totalN += tc.n
		totalQ += tc.q
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for totalN < totalNLimit && totalQ < totalQLimit {
		n := rng.Intn(8000) + 1
		if totalN+n > totalNLimit {
			n = totalNLimit - totalN
		}
		if n <= 0 {
			break
		}
		q := rng.Intn(8000) + 1
		if totalQ+q > totalQLimit {
			q = totalQLimit - totalQ
		}
		if q <= 0 {
			break
		}
		a := randPermutation(rng, n)
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			qs[i] = query{l: l, r: r}
		}
		tests = append(tests, testCase{n: n, q: q, a: a, qs: qs})
		totalN += n
		totalQ += q
	}
	return tests
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1, q: 2,
			a:  []int{1},
			qs: []query{{l: 1, r: 1}, {l: 1, r: 1}},
		},
		{
			n: 3, q: 3,
			a:  []int{1, 2, 3},
			qs: []query{{l: 1, r: 3}, {l: 2, r: 3}, {l: 1, r: 1}},
		},
		{
			n: 5, q: 4,
			a:  []int{5, 4, 3, 2, 1},
			qs: []query{{l: 1, r: 5}, {l: 2, r: 4}, {l: 3, r: 5}, {l: 1, r: 3}},
		},
		{
			n: 4, q: 3,
			a:  []int{2, 1, 4, 3},
			qs: []query{{l: 1, r: 4}, {l: 2, r: 3}, {l: 1, r: 2}},
		},
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
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.qs {
			fmt.Fprintf(&sb, "%d %d\n", qu.l, qu.r)
		}
	}
	return sb.String()
}

func singleInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.q)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, qu := range tc.qs {
		fmt.Fprintf(&sb, "%d %d\n", qu.l, qu.r)
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

func parseOutputs(out string, tests []testCase) ([]string, error) {
	fields := strings.Fields(out)
	total := 0
	for _, tc := range tests {
		total += tc.q
	}
	if len(fields) != total {
		return nil, fmt.Errorf("expected %d tokens, got %d", total, len(fields))
	}
	res := make([]string, total)
	for i, f := range fields {
		f = strings.ToUpper(f)
		if f == "YES" || f == "Y" {
			res[i] = "YES"
		} else {
			res[i] = "NO"
		}
	}
	return res, nil
}
