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

type testCase struct{ h, m int }

func generateTests() []testCase {
	var tests []testCase
	for h := 0; h < 24 && len(tests) < 100; h++ {
		for m := 0; m < 60 && len(tests) < 100; m++ {
			if h == 0 && m == 0 {
				continue
			}
			tests = append(tests, testCase{h, m})
		}
	}
	return tests
}

func expected(tc testCase) int {
	return 24*60 - (tc.h*60 + tc.m)
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
	rand.Seed(1)
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d %d\n", tc.h, tc.m)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		var got int
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
