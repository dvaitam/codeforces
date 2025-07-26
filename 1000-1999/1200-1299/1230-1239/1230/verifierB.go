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

type testCase struct {
	n int
	k int
	s string
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(2))
	tests := make([]testCase, 0, 100)
	// fixed edge cases
	tests = append(tests, testCase{n: 1, k: 0, s: "5"})
	tests = append(tests, testCase{n: 1, k: 1, s: "8"})
	for len(tests) < 100 {
		n := r.Intn(10) + 1
		k := r.Intn(n + 1)
		b := make([]byte, n)
		b[0] = byte('1' + r.Intn(9))
		for i := 1; i < n; i++ {
			b[i] = byte('0' + r.Intn(10))
		}
		tests = append(tests, testCase{n: n, k: k, s: string(b)})
	}
	return tests
}

func expected(t testCase) string {
	if t.n == 1 {
		if t.k > 0 {
			return "0"
		}
		return t.s
	}
	b := []byte(t.s)
	k := t.k
	if k > 0 && b[0] != '1' {
		b[0] = '1'
		k--
	}
	for i := 1; i < t.n && k > 0; i++ {
		if b[i] != '0' {
			b[i] = '0'
			k--
		}
	}
	return string(b)
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("run failed: %v\n%s", err, errb.String())
	}
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n%s\n", t.n, t.k, t.s)
		want := expected(t)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("test %d failed:\ninput: %q %d\nexpected %s got %s\n", i+1, t.s, t.k, want, strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
