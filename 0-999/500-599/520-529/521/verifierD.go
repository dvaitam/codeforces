package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type operation struct {
	typ int
	idx int
	val int64
}

type testCase struct {
	k, n, m int
	skills  []int64
	ops     []operation
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-521D-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", bin, "521D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.k, tc.n, tc.m))
	for i, val := range tc.skills {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(val, 10))
	}
	sb.WriteByte('\n')
	for _, op := range tc.ops {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", op.typ, op.idx, op.val))
	}
	return sb.String()
}

func parseOutput(out string, maxOps int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	l64, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number of operations: %v", err)
	}
	if l64 < 0 {
		return nil, fmt.Errorf("negative number of operations")
	}
	if l64 > int64(maxOps) {
		return nil, fmt.Errorf("uses %d operations but limit is %d", l64, maxOps)
	}
	l := int(l64)
	if len(fields) != l+1 {
		return nil, fmt.Errorf("expected %d operation indices, got %d", l, len(fields)-1)
	}
	res := make([]int, l)
	for i := 0; i < l; i++ {
		val, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid operation index %q: %v", fields[i+1], err)
		}
		res[i] = val
	}
	return res, nil
}

func evaluate(tc testCase, seq []int) (*big.Int, error) {
	if len(seq) > tc.m {
		return nil, fmt.Errorf("uses %d operations but limit is %d", len(seq), tc.m)
	}
	vals := make([]*big.Int, tc.k)
	for i := 0; i < tc.k; i++ {
		vals[i] = big.NewInt(tc.skills[i])
	}
	used := make(map[int]bool, len(seq))
	for pos, idx := range seq {
		if idx < 1 || idx > tc.n {
			return nil, fmt.Errorf("operation index %d out of range", idx)
		}
		if used[idx] {
			return nil, fmt.Errorf("operation %d repeated", idx)
		}
		used[idx] = true
		op := tc.ops[idx-1]
		if op.idx < 1 || op.idx > tc.k {
			return nil, fmt.Errorf("operation %d targets invalid skill %d", idx, op.idx)
		}
		target := vals[op.idx-1]
		switch op.typ {
		case 1:
			target.SetInt64(op.val)
		case 2:
			target.Add(target, big.NewInt(op.val))
		case 3:
			target.Mul(target, big.NewInt(op.val))
		default:
			return nil, fmt.Errorf("operation %d has invalid type %d", idx, op.typ)
		}
		if target.Sign() <= 0 {
			return nil, fmt.Errorf("skill %d became non-positive after operation %d", op.idx, idx)
		}
		_ = pos
	}
	prod := big.NewInt(1)
	for _, v := range vals {
		prod.Mul(prod, v)
	}
	return prod, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			k: 1, n: 0, m: 0,
			skills: []int64{5},
			ops:    []operation{},
		},
		{
			k: 2, n: 3, m: 2,
			skills: []int64{2, 3},
			ops: []operation{
				{typ: 1, idx: 1, val: 10},
				{typ: 2, idx: 1, val: 5},
				{typ: 3, idx: 2, val: 3},
			},
		},
		{
			k: 3, n: 4, m: 3,
			skills: []int64{1, 2, 3},
			ops: []operation{
				{typ: 2, idx: 1, val: 4},
				{typ: 1, idx: 2, val: 7},
				{typ: 3, idx: 3, val: 2},
				{typ: 2, idx: 2, val: 3},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	k := rng.Intn(6) + 1
	n := rng.Intn(20)
	m := 0
	if n > 0 {
		m = rng.Intn(n + 1)
	}
	skills := make([]int64, k)
	for i := range skills {
		if rng.Intn(4) == 0 {
			skills[i] = int64(rng.Intn(1000000) + 1)
		} else {
			skills[i] = int64(rng.Intn(20) + 1)
		}
	}
	ops := make([]operation, n)
	for i := 0; i < n; i++ {
		typ := rng.Intn(3) + 1
		idx := rng.Intn(k) + 1
		var val int64
		switch typ {
		case 1:
			if rng.Intn(4) == 0 {
				val = int64(rng.Intn(1000000) + 1)
			} else {
				val = int64(rng.Intn(30) + 1)
			}
		case 2:
			if rng.Intn(4) == 0 {
				val = int64(rng.Intn(100000) + 1)
			} else {
				val = int64(rng.Intn(20) + 1)
			}
		case 3:
			val = int64(rng.Intn(5) + 2)
		}
		ops[i] = operation{typ: typ, idx: idx, val: val}
	}
	return testCase{k: k, n: n, m: m, skills: skills, ops: ops}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expSeq, err := parseOutput(expOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}
		expProd, err := evaluate(tc, expSeq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, expOut)
			os.Exit(1)
		}

		actOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		actSeq, err := parseOutput(actOut, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, actOut, input)
			os.Exit(1)
		}
		actProd, err := evaluate(tc, actSeq)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target produced invalid sequence on test %d: %v\ninput:\n%s\noutput:\n%s\n", idx+1, err, input, actOut)
			os.Exit(1)
		}
		if actProd.Cmp(expProd) != 0 {
			fmt.Fprintf(os.Stderr, "test %d failed: expected product %s but got %s\ninput:\n%s\n", idx+1, expProd.String(), actProd.String(), input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
