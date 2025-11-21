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
	input  string
	tokens int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleF")
	cmd := exec.Command("go", "build", "-o", out, "2155F.go")
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

func parseTokens(out string) []string {
	return strings.Fields(out)
}

func deterministicTests() []testCase {
	var tests []testCase
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString("3 5 10 4\n")
	sb.WriteString("1 3\n2 1\n")
	sb.WriteString("1 1\n1 2\n1 3\n1 4\n1 5\n")
	sb.WriteString("2 1\n2 2\n2 5\n")
	sb.WriteString("3 1\n3 2\n")
	sb.WriteString("2 3\n2 1\n1 1\n")
	tests = append(tests, testCase{input: sb.String(), tokens: 4})

	var sb2 strings.Builder
	sb2.WriteString("1\n")
	sb2.WriteString("4 3 5 3\n")
	sb2.WriteString("1 2\n2 3\n3 4\n")
	sb2.WriteString("1 1\n2 1\n3 2\n4 2\n4 3\n")
	sb2.WriteString("1 3\n2 4\n1 4\n")
	tests = append(tests, testCase{input: sb2.String(), tokens: 3})

	return tests
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count)
	totalN := 0
	totalK := 0
	totalS := 0
	totalQ := 0
	for len(tests) < count && totalN < 300000 {
		t := rnd.Intn(3) + 1
		cases := make([]string, 0, t)
		tokenSum := 0
		for caseIdx := 0; caseIdx < t; caseIdx++ {
			n := rnd.Intn(50) + 1
			k := rnd.Intn(50) + 1
			s := rnd.Intn(n * k)
			if s > 200 {
				s = 200
			}
			q := rnd.Intn(50) + 1
			if totalN+n > 300000 || totalK+k > 300000 || totalS+s > 300000 || totalQ+q > 300000 {
				break
			}
			totalN += n
			totalK += k
			totalS += s
			totalQ += q
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, s, q))
			for node := 2; node <= n; node++ {
				parent := rnd.Intn(node-1) + 1
				sb.WriteString(fmt.Sprintf("%d %d\n", node, parent))
			}
			used := make(map[[2]int]bool)
			for i := 0; i < s; i++ {
				v := rnd.Intn(n) + 1
				c := rnd.Intn(k) + 1
				key := [2]int{v, c}
				if used[key] {
					i--
					continue
				}
				used[key] = true
				sb.WriteString(fmt.Sprintf("%d %d\n", v, c))
			}
			for i := 0; i < q; i++ {
				u := rnd.Intn(n) + 1
				v := rnd.Intn(n) + 1
				sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
			}
			tokenSum += q
			cases = append(cases, sb.String())
		}
		if len(cases) == 0 {
			continue
		}
		var full strings.Builder
		full.WriteString(fmt.Sprintf("%d\n", len(cases)))
		for _, block := range cases {
			full.WriteString(block)
		}
		tests = append(tests, testCase{input: full.String(), tokens: tokenSum})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
		expTokens := parseTokens(expOut)
		if len(expTokens) != tc.tokens {
			fmt.Fprintf(os.Stderr, "case %d: oracle token mismatch, expected %d got %d\n", idx+1, tc.tokens, len(expTokens))
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotTokens := parseTokens(gotOut)
		if len(gotTokens) != len(expTokens) {
			fmt.Fprintf(os.Stderr, "case %d: token count mismatch: expected %d got %d\n", idx+1, len(expTokens), len(gotTokens))
			os.Exit(1)
		}
		for i := range expTokens {
			if gotTokens[i] != expTokens[i] {
				fmt.Fprintf(os.Stderr, "case %d token %d mismatch: expected %s got %s\n", idx+1, i+1, expTokens[i], gotTokens[i])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
