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

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	outPath := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "239B.go")
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(output))
	}
	return outPath, nil
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

func buildInput(n, q int, s string, queries [][2]int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n%s\n", n, q, s)
	for _, qr := range queries {
		fmt.Fprintf(&sb, "%d %d\n", qr[0], qr[1])
	}
	return sb.String()
}

func randomSequence(n int, rnd *rand.Rand) string {
	chars := []byte("<>0123456789")
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = chars[rnd.Intn(len(chars))]
	}
	return string(res)
}

func randomQueries(n, q int, rnd *rand.Rand) [][2]int {
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rnd.Intn(n) + 1
		r := rnd.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
	}
	return queries
}

func parseOutput(out string, expectedLines int) ([][]int64, error) {
	fields := strings.Fields(out)
	required := expectedLines * 10
	if len(fields) != required {
		return nil, fmt.Errorf("expected %d integers, got %d", required, len(fields))
	}
	results := make([][]int64, expectedLines)
	for i := 0; i < expectedLines; i++ {
		results[i] = make([]int64, 10)
		for j := 0; j < 10; j++ {
			val, err := strconv.ParseInt(fields[i*10+j], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", fields[i*10+j])
			}
			results[i][j] = val
		}
	}
	return results, nil
}

func deterministicTests() []string {
	var tests []string
	tests = append(tests, buildInput(1, 1, "5", [][2]int{{1, 1}}))
	tests = append(tests, buildInput(3, 2, "<5>", [][2]int{{1, 3}, {2, 2}}))
	tests = append(tests, buildInput(5, 3, "9<8>7", [][2]int{{1, 5}, {2, 4}, {3, 5}}))
	tests = append(tests, buildInput(4, 2, "<<>>", [][2]int{{1, 4}, {2, 3}}))
	tests = append(tests, buildInput(6, 3, "123456", [][2]int{{1, 6}, {1, 3}, {4, 6}}))
	return tests
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		n := rnd.Intn(100) + 1
		q := rnd.Intn(100) + 1
		seq := randomSequence(n, rnd)
		queries := randomQueries(n, q, rnd)
		tests = append(tests, buildInput(n, q, seq, queries))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	tests = append(tests, randomTests(300)...)

	for idx, input := range tests {
		expOut, err := runProgram(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		firstLineEnd := strings.IndexByte(input, '\n')
		if firstLineEnd == -1 {
			fmt.Fprintf(os.Stderr, "case %d: malformed test input\n", idx+1)
			os.Exit(1)
		}
		var headerN, q int
		if _, err := fmt.Sscanf(input[:firstLineEnd], "%d %d", &headerN, &q); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to parse header: %v\n", idx+1, err)
			os.Exit(1)
		}
		_ = headerN

		expVals, err := parseOutput(expOut, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			for d := 0; d < 10; d++ {
				if gotVals[i][d] != expVals[i][d] {
					fmt.Fprintf(os.Stderr, "case %d, query %d digit %d mismatch: expected %d got %d\n",
						idx+1, i+1, d, expVals[i][d], gotVals[i][d])
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
