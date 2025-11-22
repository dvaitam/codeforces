package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type testCase struct {
	n    int
	grid []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Println("failed to build oracle:", err)
		return
	}
	defer cleanup()

	tests := generateTests()
	input := formatInput(tests)

	oracleOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(oracleOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, oracleOut)
		return
	}

	candOut, err := runBinary(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Failed test input:")
			fmt.Println(prettyCase(tests[i]))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2062H-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleH")
	cmd := exec.Command("go", "build", "-o", outPath, "2062H.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return outPath, cleanup, nil
}

func runBinary(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string, count int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != count {
		return nil, fmt.Errorf("expected %d outputs, got %d", count, len(fields))
	}
	res := make([]int64, count)
	for i, tok := range fields {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	sb.Grow(len(tests) * 256)
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return []byte(sb.String())
}

func prettyCase(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, grid: []string{"0"}},
		{n: 1, grid: []string{"1"}},
		{n: 2, grid: []string{"01", "10"}},
		{n: 3, grid: []string{"010", "000", "101"}},
		{n: 4, grid: []string{"0100", "1010", "0001", "1100"}},
		{n: 4, grid: []string{"1111", "1111", "1111", "1111"}},
		{n: 5, grid: []string{"10001", "01010", "00100", "01010", "10001"}},
	}
}

func generateTests() []testCase {
	tests := deterministicTests()
	const limit = 1 << 14
	used := 0
	for _, tc := range tests {
		used += 1 << tc.n
	}
	remaining := limit - used

	rng := rand.New(rand.NewSource(20622062))
	target := 90
	attempts := 0
	for len(tests) < target && remaining > 0 && attempts < 2000 {
		attempts++
		n := rng.Intn(14) + 1
		size := 1 << n
		if size > remaining {
			n = bits.Len(uint(remaining)) - 1
			if n <= 0 {
				break
			}
			size = 1 << n
		}
		onesProb := rng.Intn(61) + 20 // 20..80 for varied density
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			var sb strings.Builder
			for j := 0; j < n; j++ {
				if rng.Intn(100) < onesProb {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
			grid[i] = sb.String()
		}
		tests = append(tests, testCase{n: n, grid: grid})
		remaining -= size
	}
	return tests
}
