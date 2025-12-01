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

const refSource = "./2138A.go"

type testCase struct {
	name    string
	input   string
	k       []uint64
	x       []uint64
	answers int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()

	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refSeqs, err := parseOutput(refOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candSeqs, err := parseOutput(candOut, tc.answers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.answers; caseIdx++ {
			targetK := tc.k[caseIdx]
			targetX := tc.x[caseIdx]
			refOps := refSeqs[caseIdx]
			candOps := candSeqs[caseIdx]

			if len(candOps) != len(refOps) {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: expected minimal steps %d, got %d\n", idx+1, tc.name, caseIdx+1, len(refOps), len(candOps))
				os.Exit(1)
			}
			if err := simulate(targetK, targetX, candOps); err != nil {
				fmt.Fprintf(os.Stderr, "test %d (%s) case %d: invalid operations: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, caseIdx+1, err, tc.input, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2138A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2138A.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, cases int) ([][]int, error) {
	fields := strings.Fields(output)
	pos := 0
	res := make([][]int, cases)
	for i := 0; i < cases; i++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("missing operation count for case %d", i+1)
		}
		nOps64, err := strconv.ParseInt(fields[pos], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid operation count %q in case %d", fields[pos], i+1)
		}
		if nOps64 < 0 || nOps64 > 120 {
			return nil, fmt.Errorf("operation count out of range in case %d: %d", i+1, nOps64)
		}
		nOps := int(nOps64)
		pos++
		if pos+nOps > len(fields) {
			return nil, fmt.Errorf("not enough operations for case %d", i+1)
		}
		ops := make([]int, nOps)
		for j := 0; j < nOps; j++ {
			val, err := strconv.Atoi(fields[pos+j])
			if err != nil || (val != 1 && val != 2) {
				return nil, fmt.Errorf("invalid operation value %q in case %d", fields[pos+j], i+1)
			}
			ops[j] = val
		}
		pos += nOps
		res[i] = ops
	}
	if pos != len(fields) {
		return nil, fmt.Errorf("extra tokens at end of output")
	}
	return res, nil
}

func simulate(k, x uint64, ops []int) error {
	total := uint64(1) << (k + 1)
	a := uint64(1) << k
	b := a

	for idx, op := range ops {
		switch op {
		case 1:
			if a%2 == 0 {
				half := a / 2
				a -= half
				b += half
			} else {
				return fmt.Errorf("step %d: operation 1 when Chocola has odd cakes (%d)", idx+1, a)
			}
		case 2:
			if b%2 == 0 {
				half := b / 2
				b -= half
				a += half
			} else {
				return fmt.Errorf("step %d: operation 2 when Vanilla has odd cakes (%d)", idx+1, b)
			}
		default:
			return fmt.Errorf("invalid operation %d at step %d", op, idx+1)
		}
	}

	if a != x || b != total-x {
		return fmt.Errorf("incorrect final distribution: got (%d,%d), expected (%d,%d)", a, b, x, total-x)
	}
	return nil
}

func buildTests() []testCase {
	tests := []testCase{
		newTestCase("sample", "4\n2 3\n2 4\n3 7\n2 5\n"),
		buildSingle("already_ok", 2, 4),
		buildSingle("simple_move", 2, 6),
		buildSingle("k1", 1, 1),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomTestCase(rng, i+1))
	}
	return tests
}

func buildSingle(name string, k uint64, x uint64) testCase {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", k, x))
	return testCase{name: name, input: sb.String(), k: []uint64{k}, x: []uint64{x}, answers: 1}
}

func newTestCase(name, input string) testCase {
	kVals, xVals, cases, err := parseInputMeta(input)
	if err != nil {
		panic(fmt.Sprintf("failed to parse test %s: %v", name, err))
	}
	return testCase{name: name, input: input, k: kVals, x: xVals, answers: cases}
}

func parseInputMeta(input string) ([]uint64, []uint64, int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, nil, 0, fmt.Errorf("failed to read t: %v", err)
	}
	if t <= 0 {
		return nil, nil, 0, fmt.Errorf("non-positive t: %d", t)
	}
	kVals := make([]uint64, t)
	xVals := make([]uint64, t)
	for i := 0; i < t; i++ {
		var k, x uint64
		if _, err := fmt.Fscan(reader, &k, &x); err != nil {
			return nil, nil, 0, fmt.Errorf("failed to read case %d: %v", i+1, err)
		}
		kVals[i] = k
		xVals[i] = x
	}
	return kVals, xVals, t, nil
}

func randomTestCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	ks := make([]uint64, t)
	xs := make([]uint64, t)
	for i := 0; i < t; i++ {
		k := uint64(rng.Intn(20) + 1)
		total := uint64(1) << (k + 1)
		x := uint64(rng.Intn(int(total-1))) + 1
		ks[i] = k
		xs[i] = x
		sb.WriteString(fmt.Sprintf("%d %d\n", k, x))
	}
	return testCase{
		name:    fmt.Sprintf("random_%d", idx),
		input:   sb.String(),
		k:       ks,
		x:       xs,
		answers: t,
	}
}
