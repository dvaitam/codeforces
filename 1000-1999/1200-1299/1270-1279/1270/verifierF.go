package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

func solveRef(input string) (int64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return 0, err
	}
	n := len(s)
	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i]
		if s[i] == '1' {
			prefix[i+1]++
		}
	}
	B := int(math.Sqrt(float64(n))) + 1
	var ans int64
	for k := 1; k <= B; k++ {
		freq := make(map[int]int)
		for j := 0; j <= n; j++ {
			val := j - k*prefix[j]
			ans += int64(freq[val])
			freq[val]++
		}
	}
	maxd := n / B
	for d := 1; d <= maxd; d++ {
		limit := d * (B + 1)
		for r := 0; r < d; r++ {
			freq := make(map[int]int)
			i := r
			for j := r; j <= n; j += d {
				for i <= j-limit {
					freq[prefix[i]]++
					i += d
				}
				ans += int64(freq[prefix[j]-d])
			}
		}
	}
	return ans, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	if stderr.Len() > 0 {
		return "", fmt.Errorf("unexpected stderr output: %s", stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string) (int64, error) {
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	var val int64
	if _, err := fmt.Sscan(out, &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val, nil
}

func makeCase(name, s string) testCase {
	return testCase{name: name, input: fmt.Sprintf("%s\n", s)}
}

func randomString(rng *rand.Rand, n int) string {
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single1", "1"),
		makeCase("single0", "0"),
		makeCase("allones", "11111"),
		makeCase("alternating", "101010"),
		makeCase("prefixzeros", "00001111"),
		makeCase("mixed", "101001101"),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 60; i++ {
		n := rng.Intn(200) + 1
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", i+1), randomString(rng, n)))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out)
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%d\nactual:%d\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
