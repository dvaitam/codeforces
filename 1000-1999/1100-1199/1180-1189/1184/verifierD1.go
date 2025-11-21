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

type testCase struct {
	input string
	t     int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleD1")
	cmd := exec.Command("go", "build", "-o", out, "1184D1.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return out, nil
}

func runProgram(bin string, input string) (string, error) {
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

func parseOutput(out string, t int) ([][2]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t*2 {
		return nil, fmt.Errorf("expected %d integers (for %d steps), got %d", t*2, t, len(fields))
	}
	ans := make([][2]int, t)
	for i := 0; i < t; i++ {
		lVal, err := strconv.Atoi(fields[2*i])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[2*i])
		}
		kVal, err := strconv.Atoi(fields[2*i+1])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[2*i+1])
		}
		ans[i] = [2]int{lVal, kVal}
	}
	return ans, nil
}

func deterministicTests() []testCase {
	tests := []testCase{}
	// Sample from problem note
	tests = append(tests, makeTest(5, 2, 5, [][2]int{
		{0, 1},
		{1, 1},
		{0, 4},
		{1, 2},
	}))
	// Only insertions
	tests = append(tests, makeTest(3, 2, 6, [][2]int{
		{1, 1},
		{1, 4},
		{1, 3},
	}))
	// Only cuts
	tests = append(tests, makeTest(6, 4, 6, [][2]int{
		{0, 3},
		{0, 2},
		{0, 1},
	}))
	return tests
}

func makeTest(n, k, m int, ops [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, m, len(ops)))
	for _, op := range ops {
		sb.WriteString(fmt.Sprintf("%d %d\n", op[0], op[1]))
	}
	return testCase{input: sb.String(), t: len(ops)}
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rnd.Intn(249) + 2   // 2..250
		m := rnd.Intn(251-n) + n // n..250
		t := rnd.Intn(1000) + 1  // 1..1000
		k := rnd.Intn(n-1) + 1   // 1..n-1? but need inclusive n
		if n == 1 {
			k = 1
		} else {
			k = rnd.Intn(n) + 1
		}
		length := n
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, m, t))
		for step := 0; step < t; step++ {
			var typ int
			switch {
			case length == 1:
				typ = 1
			case length == m:
				typ = 0
			default:
				typ = rnd.Intn(2)
			}
			if typ == 1 {
				pos := rnd.Intn(length+1) + 1
				sb.WriteString(fmt.Sprintf("1 %d\n", pos))
				if pos <= k {
					k++
				}
				length++
			} else {
				// ensure length>=2
				if length < 2 {
					pos := 1
					sb.WriteString(fmt.Sprintf("1 %d\n", pos))
					if pos <= k {
						k++
					}
					length++
					continue
				}
				pos := rnd.Intn(length-1) + 1
				sb.WriteString(fmt.Sprintf("0 %d\n", pos))
				if k <= pos {
					length = pos
				} else {
					length = length - pos
					k = k - pos
				}
			}
		}
		tests = append(tests, testCase{input: sb.String(), t: t})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(200)...)

	for idx, tc := range tests {
		expOut, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		expPairs, err := parseOutput(expOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid oracle output: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotPairs, err := parseOutput(gotOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		for step := 0; step < tc.t; step++ {
			if gotPairs[step] != expPairs[step] {
				fmt.Fprintf(os.Stderr, "case %d step %d mismatch: expected (%d %d) got (%d %d)\n",
					idx+1, step+1, expPairs[step][0], expPairs[step][1], gotPairs[step][0], gotPairs[step][1])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
