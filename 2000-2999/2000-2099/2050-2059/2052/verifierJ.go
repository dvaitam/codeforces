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
	"time"
)

const (
	refSource  = "2000-2999/2000-2099/2050-2059/2052/2052J.go"
	totalLimit = 50000
)

type testCase struct {
	n int
	m int
	q int

	a []int64
	d []int64
	l []int64
	t []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

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
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i := range tests {
		for j := 0; j < tests[i].q; j++ {
			if candAns[i][j] != refAns[i][j] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d query %d: expected %d, got %d\n", i+1, j+1, refAns[i][j], candAns[i][j])
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2052J-ref-*")
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d\n", tc.n, tc.m, tc.q)
		for i, v := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for i, v := range tc.d {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for i, v := range tc.l {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
		for i, v := range tc.t {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func parseOutput(out string, tests []testCase) ([][]int, error) {
	fields := strings.Fields(out)
	totalNeed := 0
	for _, tc := range tests {
		totalNeed += tc.q
	}
	if len(fields) < totalNeed {
		return nil, fmt.Errorf("not enough numbers: expected %d, got %d", totalNeed, len(fields))
	}

	res := make([][]int, len(tests))
	idx := 0
	for i, tc := range tests {
		res[i] = make([]int, tc.q)
		for j := 0; j < tc.q; j++ {
			val, err := strconv.ParseInt(fields[idx], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse integer %q at position %d: %v", fields[idx], idx+1, err)
			}
			res[i][j] = int(val)
			if res[i][j] < 0 || res[i][j] > tc.m {
				return nil, fmt.Errorf("value out of range on test %d query %d: %d", i+1, j+1, res[i][j])
			}
			idx++
		}
	}
	if idx != len(fields) {
		return nil, fmt.Errorf("extra data after reading outputs (consumed %d of %d)", idx, len(fields))
	}
	return res, nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	total := 0

	add := func(tc testCase) {
		need := tc.n + tc.m + tc.q
		if total+need > totalLimit {
			return
		}
		tests = append(tests, tc)
		total += need
	}

	// Simple edge cases.
	add(simpleCase([]int64{10}, []int64{10}, []int64{5}, []int64{5, 10, 11}))
	add(simpleCase([]int64{3, 4}, []int64{3, 7}, []int64{2, 2, 2}, []int64{2, 6, 20, 1}))

	for attempts := 0; attempts < 400 && total < totalLimit; attempts++ {
		n := rng.Intn(25) + 1
		m := rng.Intn(25) + 1
		q := rng.Intn(25) + 1

		a, d := buildFeasibleTasks(n, rng)
		l := make([]int64, m)
		for i := range l {
			if rng.Intn(8) == 0 {
				l[i] = 1_000_000_000
			} else {
				l[i] = int64(rng.Intn(1_000_000) + 1)
			}
		}
		lastDeadline := maxInt64Slice(d)
		if lastDeadline < 1 {
			lastDeadline = 1
		}
		times := make([]int64, q)
		for i := 0; i < q; i++ {
			switch rng.Intn(4) {
			case 0:
				times[i] = int64(rng.Intn(10)+1) + lastDeadline/4
			case 1:
				times[i] = lastDeadline - int64(rng.Intn(5))*(maxInt64(1, lastDeadline/10))
			case 2:
				times[i] = lastDeadline + int64(rng.Intn(500_000)+1)
			default:
				times[i] = lastDeadline/2 + int64(rng.Intn(1_000_000))
			}
			if times[i] < 1 {
				times[i] = 1
			}
			if times[i] > 1_000_000_000_000_000 {
				times[i] = 1_000_000_000_000_000
			}
		}

		add(testCase{n: n, m: m, q: q, a: a, d: d, l: l, t: times})
	}

	if len(tests) == 0 {
		tests = append(tests, simpleCase([]int64{1}, []int64{1}, []int64{1}, []int64{1}))
	}
	return tests
}

func simpleCase(a, d, l, t []int64) testCase {
	return testCase{
		n: len(a),
		m: len(l),
		q: len(t),
		a: append([]int64(nil), a...),
		d: append([]int64(nil), d...),
		l: append([]int64(nil), l...),
		t: append([]int64(nil), t...),
	}
}

// buildFeasibleTasks returns task durations and deadlines that are guaranteed to be feasible.
func buildFeasibleTasks(n int, rng *rand.Rand) ([]int64, []int64) {
	a := make([]int64, n)
	pairs := make([]struct {
		a int64
		d int64
	}, n)

	var pref int64
	base := rng.Int63n(900_000_000_000_000)
	for i := 0; i < n; i++ {
		if rng.Intn(6) == 0 {
			pairs[i].a = 1_000_000_000
		} else {
			pairs[i].a = int64(rng.Intn(1_000_000) + 1)
		}
		pref += pairs[i].a
		var slack int64
		if rng.Intn(5) == 0 {
			slack = 0
		} else {
			slack = int64(rng.Intn(1_000_000))
		}
		pairs[i].d = pref + base + slack
		if pairs[i].d < 1 {
			pairs[i].d = 1
		}
		if pairs[i].d > 1_000_000_000_000_000 {
			pairs[i].d = 1_000_000_000_000_000
		}
	}

	rng.Shuffle(n, func(i, j int) {
		pairs[i], pairs[j] = pairs[j], pairs[i]
	})
	for i := 0; i < n; i++ {
		a[i] = pairs[i].a
	}
	d := make([]int64, n)
	for i := 0; i < n; i++ {
		d[i] = pairs[i].d
	}
	return a, d
}

func maxInt64Slice(arr []int64) int64 {
	var best int64
	for i, v := range arr {
		if i == 0 || v > best {
			best = v
		}
	}
	return best
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
