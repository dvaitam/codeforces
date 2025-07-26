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

type testCase struct{ n, k int64 }

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(1))
	tests := make([]testCase, 100)
	for i := range tests {
		tests[i].n = rnd.Int63n(1e9) + 1
		tests[i].k = rnd.Int63n(1e9) + 1
	}
	return tests
}

func expected(tc testCase) int64 {
	base := (tc.n / tc.k) * tc.k
	rem := tc.n % tc.k
	half := tc.k / 2
	if rem > half {
		rem = half
	}
	return base + rem
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), fmt.Errorf("exec error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.n, tc.k)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Printf("test %d: failed to parse output '%s'\n", i+1, strings.TrimSpace(out))
			os.Exit(1)
		}
		exp := expected(tc)
		if got != exp {
			fmt.Printf("test %d: expected %d got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
