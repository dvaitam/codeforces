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

type testCase struct {
	m      int
	d      int
	spaces [][][]int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-958D2-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleD2")
	cmd := exec.Command("go", "build", "-o", path, "958D2.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return path, cleanup, nil
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
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.m, tc.d))
	for _, space := range tc.spaces {
		sb.WriteString(fmt.Sprintf("%d\n", len(space)))
		for _, vec := range space {
			for j, val := range vec {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(val))
			}
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseAnswer(out string, m int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != m {
		return nil, fmt.Errorf("expected %d numbers, got %d", m, len(fields))
	}
	res := make([]int, m)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		if val < 1 || val > m {
			return nil, fmt.Errorf("group id %d out of range", val)
		}
		res[i] = val
	}
	return res, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			m: 2, d: 1,
			spaces: [][][]int{
				{{1}},
				{{2}},
			},
		},
		{
			m: 3, d: 2,
			spaces: [][][]int{
				{{1, 0}},
				{{2, 0}, {4, 0}},
				{{0, 1}},
			},
		},
		{
			m: 4, d: 3,
			spaces: [][][]int{
				{{1, 0, 0}, {0, 1, 0}},
				{{0, 1, 0}, {1, 0, 0}},
				{{0, 0, 1}},
				{{0, 0, 1}, {0, 0, 2}},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	d := rng.Intn(5) + 1
	m := rng.Intn(50) + 1
	spaces := make([][][]int, m)
	for i := 0; i < m; i++ {
		k := rng.Intn(d) + 1
		space := make([][]int, k)
		for j := 0; j < k; j++ {
			vec := make([]int, d)
			for t := 0; t < d; t++ {
				vec[t] = rng.Intn(501) - 250
			}
			space[j] = vec
		}
		spaces[i] = space
	}
	return testCase{m: m, d: d, spaces: spaces}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
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

		expStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expStr, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expStr)
			os.Exit(1)
		}

		gotStr, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotStr, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotStr, input)
			os.Exit(1)
		}

		if !equalSlices(exp, got) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d:\nexpected: %v\ngot: %v\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
