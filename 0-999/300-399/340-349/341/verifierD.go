package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	refSource        = "341D.go"
	tempOraclePrefix = "oracle-341D-"
	randomTests      = 120
	maxRandomN       = 40
	maxRandomOps     = 200
)

type operation struct {
	typ            int
	x0, y0, x1, y1 int
	v              uint64
}

type testCase struct {
	name       string
	n          int
	operations []operation
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTestsCases(randomTests, rng)...)

	for idx, tc := range tests {
		input, queryCount := buildInput(tc)
		expOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		expected, err := parseOutputs(expOut, queryCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := parseOutputs(candOut, queryCount)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, candOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		for i := range expected {
			if expected[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at query %d: expected %d got %d\n", idx+1, tc.name, i+1, expected[i], got[i])
				fmt.Println("Input:")
				fmt.Print(input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testCase) (string, int) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, len(tc.operations))
	queryCount := 0
	for _, op := range tc.operations {
		if op.typ == 1 {
			fmt.Fprintf(&sb, "1 %d %d %d %d\n", op.x0, op.y0, op.x1, op.y1)
			queryCount++
		} else {
			fmt.Fprintf(&sb, "2 %d %d %d %d %d\n", op.x0, op.y0, op.x1, op.y1, op.v)
		}
	}
	return sb.String(), queryCount
}

func parseOutputs(out string, expected int) ([]uint64, error) {
	if expected == 0 {
		if strings.TrimSpace(out) == "" {
			return nil, nil
		}
		return nil, fmt.Errorf("unexpected output for input without queries")
	}
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(fields))
	}
	res := make([]uint64, expected)
	for i, f := range fields {
		val, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not an unsigned integer", f)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	ops1 := []operation{
		{typ: 1, x0: 1, y0: 1, x1: 1, y1: 1},
		{typ: 2, x0: 1, y0: 1, x1: 1, y1: 1, v: 5},
		{typ: 1, x0: 1, y0: 1, x1: 1, y1: 1},
		{typ: 2, x0: 1, y0: 1, x1: 1, y1: 1, v: 5},
		{typ: 1, x0: 1, y0: 1, x1: 1, y1: 1},
	}
	ops2 := []operation{
		{typ: 2, x0: 1, y0: 1, x1: 2, y1: 2, v: 3},
		{typ: 2, x0: 2, y0: 2, x1: 3, y1: 3, v: 7},
		{typ: 1, x0: 1, y0: 1, x1: 3, y1: 3},
		{typ: 1, x0: 2, y0: 2, x1: 2, y1: 3},
		{typ: 2, x0: 1, y0: 1, x1: 3, y1: 3, v: 7},
		{typ: 1, x0: 1, y0: 1, x1: 3, y1: 3},
	}
	ops3 := []operation{
		{typ: 2, x0: 1, y0: 1, x1: 4, y1: 4, v: 1},
		{typ: 2, x0: 2, y0: 2, x1: 4, y1: 4, v: 2},
		{typ: 2, x0: 3, y0: 3, x1: 4, y1: 4, v: 4},
		{typ: 1, x0: 1, y0: 1, x1: 4, y1: 4},
		{typ: 1, x0: 2, y0: 2, x1: 3, y1: 3},
		{typ: 1, x0: 4, y0: 4, x1: 4, y1: 4},
	}
	ops4 := []operation{
		{typ: 1, x0: 1, y0: 1, x1: 5, y1: 5},
	}
	ops5 := []operation{
		{typ: 2, x0: 1, y0: 1, x1: 1000, y1: 1000, v: 1},
		{typ: 1, x0: 1, y0: 1, x1: 1000, y1: 1000},
		{typ: 2, x0: 500, y0: 500, x1: 1000, y1: 1000, v: 2},
		{typ: 1, x0: 400, y0: 400, x1: 600, y1: 600},
	}
	return []testCase{
		{name: "single_cell_toggle", n: 1, operations: ops1},
		{name: "overlapping_rects", n: 3, operations: ops2},
		{name: "nested_updates", n: 4, operations: ops3},
		{name: "query_only", n: 5, operations: ops4},
		{name: "large_n_edges", n: 1000, operations: ops5},
	}
}

func randomTestsCases(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxRandomN) + 1
		opCount := rng.Intn(maxRandomOps) + 1
		ops := make([]operation, 0, opCount)
		for j := 0; j < opCount; j++ {
			typ := 2
			if rng.Intn(3) == 0 {
				typ = 1
			}
			if len(ops) == 0 && typ == 1 {
				typ = 2
			}
			x0, x1 := randomRange(n, rng)
			y0, y1 := randomRange(n, rng)
			if typ == 1 {
				ops = append(ops, operation{typ: 1, x0: x0, y0: y0, x1: x1, y1: y1})
			} else {
				val := uint64(rng.Int63()) & ((uint64(1) << 62) - 1)
				ops = append(ops, operation{typ: 2, x0: x0, y0: y0, x1: x1, y1: y1, v: val})
			}
		}
		if !haveQuery(ops) {
			// ensure at least one query
			x0, x1 := randomRange(n, rng)
			y0, y1 := randomRange(n, rng)
			ops = append(ops, operation{typ: 1, x0: x0, y0: y0, x1: x1, y1: y1})
		}
		tests = append(tests, testCase{
			name:       fmt.Sprintf("random_%d", i+1),
			n:          n,
			operations: ops,
		})
	}
	return tests
}

func haveQuery(ops []operation) bool {
	for _, op := range ops {
		if op.typ == 1 {
			return true
		}
	}
	return false
}

func randomRange(limit int, rng *rand.Rand) (int, int) {
	a := rng.Intn(limit) + 1
	b := rng.Intn(limit) + 1
	if a > b {
		a, b = b, a
	}
	return a, b
}
