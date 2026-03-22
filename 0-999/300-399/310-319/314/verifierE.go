package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded solver for 314E
func solve314E(input string) (string, error) {
	fields := strings.Fields(input)
	if len(fields) < 2 {
		return "", fmt.Errorf("bad input")
	}
	n, _ := strconv.Atoi(fields[0])
	s := fields[1]

	if n%2 != 0 {
		return "0", nil
	}

	dp := make([]uint32, n/2+2)
	dp[0] = 1

	for i := 1; i <= n; i++ {
		isQ := (s[i-1] == '?')
		minJ := i % 2
		maxJ := i
		if n-i < maxJ {
			maxJ = n - i
		}

		j := minJ
		if j == 0 {
			if isQ {
				dp[0] = dp[1]
			} else {
				dp[0] = 0
			}
			j += 2
		}

		if isQ {
			for ; j <= maxJ; j += 2 {
				dp[j] = dp[j-1]*25 + dp[j+1]
			}
		} else {
			for ; j <= maxJ; j += 2 {
				dp[j] = dp[j-1]
			}
		}
	}

	return fmt.Sprintf("%d", dp[0]), nil
}

func buildInput(n int, s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 32)
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	sb.WriteString(s)
	sb.WriteByte('\n')
	return sb.String()
}

func deterministicTests() []string {
	var tests []string
	tests = append(tests, buildInput(1, "a"))
	tests = append(tests, buildInput(2, "aa"))
	tests = append(tests, buildInput(2, "a?"))
	tests = append(tests, buildInput(2, "??"))
	tests = append(tests, buildInput(4, "abba"))
	tests = append(tests, buildInput(4, "a??b"))
	tests = append(tests, buildInput(6, "abcdef"))
	tests = append(tests, buildInput(6, "a?b?c?"))
	tests = append(tests, buildInput(8, "?a?b?c?d"))
	tests = append(tests, buildInput(100000, strings.Repeat("?", 100000)))
	return tests
}

func randomString(n int, rnd *rand.Rand) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if rnd.Intn(3) == 0 {
			buf[i] = '?'
		} else {
			c := byte('a' + rnd.Intn(26))
			if c == 'x' {
				c = 'y'
			}
			buf[i] = c
		}
	}
	return string(buf)
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%20 == 0:
			n = rnd.Intn(100000) + 1
		case i%5 == 0:
			n = rnd.Intn(1000) + 1
		default:
			n = rnd.Intn(50) + 1
		}
		tests = append(tests, buildInput(n, randomString(n, rnd)))
	}
	return tests
}

func parseOutput(out string) (uint64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		if len(fields) == 0 {
			return 0, fmt.Errorf("empty output")
		}
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := deterministicTests()
	tests = append(tests, randomTests(300)...)

	for idx, input := range tests {
		expOut, err := solve314E(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
