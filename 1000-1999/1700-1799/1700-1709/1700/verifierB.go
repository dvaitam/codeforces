package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const numTestsB = 100

func generateTestsB() []string {
	rng := rand.New(rand.NewSource(2))
	tests := make([]string, numTestsB)
	for i := 0; i < numTestsB; i++ {
		n := rng.Intn(18) + 2
		digits := make([]byte, n)
		digits[0] = byte(rng.Intn(9)+1) + '0'
		for j := 1; j < n; j++ {
			digits[j] = byte(rng.Intn(10)) + '0'
		}
		s := string(digits)
		tests[i] = fmt.Sprintf("1\n%d\n%s\n", n, s)
	}
	return tests
}

func solveB(input string) string {
	var t, n int
	var s string
	fmt.Sscan(input, &t, &n, &s)
	v := make([]int, n)
	for i := 0; i < n; i++ {
		v[i] = int(s[i] - '0')
	}
	for i := 0; i < n; i++ {
		v[i] = 9 - v[i]
	}
	if v[0] == 0 {
		carry := 0
		d := v[n-1] + 2
		carry = d / 10
		v[n-1] = d % 10
		for i := n - 2; i >= 0; i-- {
			d = v[i] + 3 + carry
			carry = d / 10
			v[i] = d % 10
		}
	}
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteByte(byte(v[i] + '0'))
	}
	b.WriteByte('\n')
	return b.String()
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsB()
	for i, tc := range tests {
		expected := solveB(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("Test %d: error running binary: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, tc, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
