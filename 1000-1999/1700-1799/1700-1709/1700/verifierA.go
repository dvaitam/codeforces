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

const numTestsA = 100

func generateTestsA() []string {
	rng := rand.New(rand.NewSource(1))
	tests := make([]string, numTestsA)
	for i := 0; i < numTestsA; i++ {
		n := rng.Intn(1000) + 1
		m := rng.Intn(1000) + 1
		tests[i] = fmt.Sprintf("1\n%d %d\n", n, m)
	}
	return tests
}

func solveA(input string) string {
	var t int
	var n, m int64
	fmt.Sscan(input, &t, &n, &m)
	ans := m*(m+1)/2 + m*(n*(n+1)/2-1)
	return fmt.Sprintf("%d\n", ans)
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
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsA()
	for i, tc := range tests {
		expected := solveA(tc)
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
