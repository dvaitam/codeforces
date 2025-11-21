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

type item struct {
	ini int
	out int
	w   int
	s   int
	v   int
}

type testCase struct {
	n  int
	S  int
	it []item
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-480D-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", bin, "480D.go")
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
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.S))
	for _, it := range tc.it {
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", it.ini, it.out, it.w, it.s, it.v))
	}
	return sb.String()
}

func parseAnswer(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected one integer, got %d: %q", len(fields), out)
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1,
			S: 0,
			it: []item{
				{ini: 0, out: 1, w: 0, s: 0, v: 5},
			},
		},
		{
			n: 2,
			S: 5,
			it: []item{
				{0, 2, 2, 3, 10},
				{1, 3, 3, 5, 7},
			},
		},
		{
			n: 3,
			S: 10,
			it: []item{
				{0, 2, 3, 7, 5},
				{1, 4, 4, 6, 6},
				{2, 5, 1, 8, 4},
			},
		},
		{
			n: 4,
			S: 7,
			it: []item{
				{0, 2, 3, 3, 4},
				{1, 3, 2, 4, 5},
				{2, 4, 1, 2, 2},
				{3, 5, 2, 3, 1},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(40) + 1
	if rng.Intn(5) == 0 {
		n = rng.Intn(100) + 1
	}
	S := rng.Intn(1000) + 1
	items := make([]item, n)
	maxTime := 2 * n
	for i := 0; i < n; i++ {
		ini := rng.Intn(maxTime - 1)
		out := ini + 1 + rng.Intn(maxTime-ini-1)
		if out >= maxTime {
			out = maxTime - 1
		}
		if out <= ini {
			out = ini + 1
		}
		w := rng.Intn(1000)
		s := rng.Intn(1000)
		v := rng.Intn(1000000) + 1
		items[i] = item{ini, out, w, s, v}
	}
	return testCase{n: n, S: S, it: items}
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
	for i := 0; i < 300; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, expStr)
			os.Exit(1)
		}

		gotStr, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotStr, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
