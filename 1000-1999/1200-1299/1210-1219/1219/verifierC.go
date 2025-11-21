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

func addOne(s string) string {
	b := []byte(s)
	i := len(b) - 1
	carry := byte(1)
	for i >= 0 && carry > 0 {
		sum := b[i] - '0' + carry
		b[i] = sum%10 + '0'
		carry = sum / 10
		i--
	}
	if carry > 0 {
		b = append([]byte{'1'}, b...)
	}
	return string(b)
}

func solveRef(input string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var L int
	if _, err := fmt.Fscan(reader, &L); err != nil {
		return "", err
	}
	var A string
	if _, err := fmt.Fscan(reader, &A); err != nil {
		return "", err
	}
	n := len(A)
	base := "1" + strings.Repeat("0", L-1)
	if n < L {
		return base, nil
	}
	if n%L != 0 {
		k := (n + L - 1) / L
		return strings.Repeat(base, k), nil
	}
	k := n / L
	prefix := A[:L]
	candidate := strings.Repeat(prefix, k)
	if candidate > A {
		return candidate, nil
	}
	inc := addOne(prefix)
	if len(inc) > L {
		return strings.Repeat(base, k+1), nil
	}
	return strings.Repeat(inc, k), nil
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

func makeCase(name string, L int, A string) testCase {
	return testCase{
		name:  name,
		input: fmt.Sprintf("%d\n%s\n", L, A),
	}
}

func randomNumber(rng *rand.Rand, length int) string {
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		if i == 0 {
			sb.WriteByte(byte(rng.Intn(9)+1) + '0')
		} else {
			sb.WriteByte(byte(rng.Intn(10)) + '0')
		}
	}
	return sb.String()
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	for i := 0; i < 120; i++ {
		L := rng.Intn(20) + 1
		length := rng.Intn(60) + 1
		A := randomNumber(rng, length)
		tests = append(tests, makeCase(fmt.Sprintf("rand_%d", i+1), L, A))
	}
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("simple1", 3, "123456"),
		makeCase("simple2", 3, "12345"),
		makeCase("equal_length", 2, "9999"),
		makeCase("small", 1, "8"),
		makeCase("large_prefix", 4, "99999999"),
		makeCase("periodic", 3, "111111"),
		makeCase("non_periodic", 3, "100000"),
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimLeft(out, "0") == "" {
			out = "0"
		}
		if strings.TrimLeft(expect, "0") == "" {
			expect = "0"
		}
		if out != expect {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%s\nactual:%s\n", idx+1, tc.name, tc.input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
