package main

import (
	"bufio"
	"bytes"
	"fmt"
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

type fraction struct {
	num int
	den int
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func parseExpression(expr string) (fraction, error) {
	var a, b, c int
	if _, err := fmt.Sscanf(expr, "(%d+%d)/%d", &a, &b, &c); err != nil {
		return fraction{}, fmt.Errorf("failed to parse %q: %v", expr, err)
	}
	n := a + b
	g := gcd(n, c)
	return fraction{num: n / g, den: c / g}, nil
}

func runReference(input string) ([]fraction, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return nil, err
	}
	result := make([]fraction, m)
	for i := 0; i < m; i++ {
		var expr string
		if _, err := fmt.Fscan(reader, &expr); err != nil {
			return nil, err
		}
		f, err := parseExpression(expr)
		if err != nil {
			return nil, err
		}
		result[i] = f
	}
	return result, nil
}

func expectedCounts(fracs []fraction) []int {
	counts := make(map[fraction]int)
	for _, f := range fracs {
		counts[f]++
	}
	res := make([]int, len(fracs))
	for i, f := range fracs {
		res[i] = counts[f]
	}
	return res
}

func makeCase(name string, exprs []string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(exprs))
	for _, e := range exprs {
		sb.WriteString(e)
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomExpression(rng *rand.Rand) string {
	a := rng.Intn(99) + 1
	b := rng.Intn(99) + 1
	c := rng.Intn(99) + 1
	return fmt.Sprintf("(%d+%d)/%d", a, b, c)
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single", []string{"(1+1)/2"}),
		makeCase("all_equal", []string{"(1+1)/2", "(2+0)/2", "(3+1)/4"}),
		makeCase("distinct", []string{"(1+2)/3", "(3+4)/5", "(7+8)/9"}),
		makeCase("mixed", []string{"(1+2)/3", "(2+4)/6", "(5+1)/3", "(3+3)/6"}),
		makeCase("repeats", []string{"(10+10)/5", "(15+5)/5", "(3+2)/1", "(7+3)/1"}),
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 60; i++ {
		n := rng.Intn(8) + 1
		exprs := make([]string, n)
		for j := 0; j < n; j++ {
			exprs[j] = randomExpression(rng)
		}
		tests = append(tests, makeCase(fmt.Sprintf("small_%d", i+1), exprs))
	}
	for i := 0; i < 40; i++ {
		n := rng.Intn(200) + 50
		exprs := make([]string, n)
		for j := 0; j < n; j++ {
			exprs[j] = randomExpression(rng)
		}
		tests = append(tests, makeCase(fmt.Sprintf("medium_%d", i+1), exprs))
	}
	for i := 0; i < 10; i++ {
		n := rng.Intn(1000) + 500
		exprs := make([]string, n)
		for j := 0; j < n; j++ {
			exprs[j] = randomExpression(rng)
		}
		tests = append(tests, makeCase(fmt.Sprintf("large_%d", i+1), exprs))
	}
	return tests
}

func runCandidate(bin string, input string) (string, error) {
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

func parseOutput(out string, m int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != m {
		return nil, fmt.Errorf("expected %d numbers, got %d", m, len(fields))
	}
	res := make([]int, m)
	for i, f := range fields {
		if _, err := fmt.Sscan(f, &res[i]); err != nil {
			return nil, fmt.Errorf("failed to parse integer #%d: %v", i+1, err)
		}
	}
	return res, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		fracs, err := runReference(tc.input)
		if err != nil {
			fmt.Printf("failed to parse test %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		expect := expectedCounts(fracs)
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(out, len(expect))
		if err != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, err, tc.input, out)
			os.Exit(1)
		}
		for i := range expect {
			if expect[i] != got[i] {
				fmt.Printf("test %d (%s) mismatch at index %d: expect %d got %d\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, i+1, expect[i], got[i], tc.input, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
