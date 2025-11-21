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

const mod = 998244353

type op struct {
	o int
	x int
}

type caseData struct {
	n, q int
	a    []int64
	b    []int64
	ops  []op
}

type testInput struct {
	text          string
	answerCounts  []int
	caseSummaries []caseData
}

func buildReference() (string, error) {
	refDir := filepath.Join("2000-2999", "2000-2099", "2050-2059", "2053")
	tmp, err := os.CreateTemp("", "ref2053D")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	tmp.Close()
	os.Remove(tmpPath)

	cmd := exec.Command("go", "build", "-o", tmpPath, "2053D.go")
	cmd.Dir = refDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
	}
	return tmpPath, nil
}

func commandForPath(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	case ".js":
		return exec.Command("node", path)
	default:
		return exec.Command(path)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := commandForPath(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func parseOutputs(output string, counts []int) ([]int64, error) {
	fields := strings.Fields(output)
	total := 0
	for _, c := range counts {
		total += c
	}
	if len(fields) != total {
		return nil, fmt.Errorf("expected %d integers in output, got %d", total, len(fields))
	}
	res := make([]int64, total)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		val %= mod
		if val < 0 {
			val += mod
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(cases []caseData) testInput {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	answerCounts := make([]int, len(cases))
	for idx, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.q))
		for i, v := range cs.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range cs.b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, op := range cs.ops {
			sb.WriteString(fmt.Sprintf("%d %d\n", op.o, op.x))
		}
		answerCounts[idx] = cs.q + 1
	}
	return testInput{
		text:          sb.String(),
		answerCounts:  answerCounts,
		caseSummaries: cases,
	}
}

func fixedInputs() []testInput {
	// Sample from statement
	case1 := caseData{
		n: 3,
		q: 4,
		a: []int64{1, 1, 2},
		b: []int64{3, 2, 1},
		ops: []op{
			{1, 3},
			{2, 3},
			{1, 1},
			{2, 3},
		},
	}
	// Construct a medium deterministic test.
	n2, q2 := 5, 5
	a2 := []int64{1, 4, 2, 7, 3}
	b2 := []int64{5, 6, 5, 6, 3}
	ops2 := []op{{2, 5}, {1, 2}, {2, 3}, {1, 1}, {2, 4}}

	// Large structured test (limits but manageable)
	nLarge := 20000
	qLarge := 20000
	aLarge := make([]int64, nLarge)
	bLarge := make([]int64, nLarge)
	opsLarge := make([]op, qLarge)
	for i := 0; i < nLarge; i++ {
		aLarge[i] = int64(1 + i%1000)
		bLarge[i] = int64(2 + (i*3)%1000)
	}
	for i := 0; i < qLarge; i++ {
		if i%2 == 0 {
			opsLarge[i] = op{1, (i % nLarge) + 1}
		} else {
			opsLarge[i] = op{2, (nLarge - i%nLarge)}
		}
	}

	return []testInput{
		buildInput([]caseData{case1}),
		buildInput([]caseData{{n: n2, q: q2, a: a2, b: b2, ops: ops2}}),
		buildInput([]caseData{{n: nLarge, q: qLarge, a: aLarge, b: bLarge, ops: opsLarge}}),
	}
}

func randomCase(rng *rand.Rand, maxN, maxQ int) caseData {
	n := rng.Intn(maxN) + 1
	q := rng.Intn(maxQ) + 1
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(500_000_000) + 1
		b[i] = rng.Int63n(500_000_000) + 1
	}
	ops := make([]op, q)
	for i := 0; i < q; i++ {
		o := 1
		if rng.Intn(2) == 1 {
			o = 2
		}
		x := rng.Intn(n) + 1
		ops[i] = op{o, x}
	}
	return caseData{n: n, q: q, a: a, b: b, ops: ops}
}

func randomInputs() []testInput {
	tests := fixedInputs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 40 {
		numCases := rng.Intn(3) + 1
		cases := make([]caseData, numCases)
		for i := 0; i < numCases; i++ {
			cases[i] = randomCase(rng, 300, 300)
		}
		tests = append(tests, buildInput(cases))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := randomInputs()
	for idx, input := range tests {
		expectOut, err := runBinary(ref, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput excerpt:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		expectVals, err := parseOutputs(expectOut, input.answerCounts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, expectOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(bin, input.text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput excerpt:\n%s\n", idx+1, err, preview(input.text))
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, input.answerCounts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for i := range expectVals {
			if expectVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d at answer %d: expected %d, got %d\ninput excerpt:\n%s\n", idx+1, i+1, expectVals[i], gotVals[i], preview(input.text))
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func preview(s string) string {
	if len(s) <= 500 {
		return s
	}
	return s[:500] + "...\n"
}
