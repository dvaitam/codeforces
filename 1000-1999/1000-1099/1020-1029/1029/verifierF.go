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

type TestCase struct {
	Input  string
	Output string
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func divisors(x int64) []int64 {
	var m []int64
	for i := int64(1); i*i <= x; i++ {
		if x%i == 0 {
			m = append(m, i)
		}
	}
	return m
}

func solveF(a, b int64) string {
	c := a + b
	am := divisors(a)
	bm := divisors(b)
	cm := divisors(c)
	for i := len(cm) - 1; i >= 0; i-- {
		d := cm[i]
		dd := c / d
		for _, da := range am {
			if d >= da && dd >= a/da {
				return fmt.Sprintf("%d", (d+dd)*2)
			}
		}
		for _, db := range bm {
			if d >= db && dd >= b/db {
				return fmt.Sprintf("%d", (d+dd)*2)
			}
		}
	}
	return ""
}

func generateTests() []TestCase {
	rand.Seed(47)
	tests := make([]TestCase, 100)
	for t := 0; t < 100; t++ {
		a := int64(rand.Intn(10000) + 1)
		b := int64(rand.Intn(10000) + 1)
		input := fmt.Sprintf("%d %d\n", a, b)
		output := solveF(a, b)
		tests[t] = TestCase{Input: input, Output: output}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.Input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(tc.Output) {
			fmt.Fprintf(os.Stderr, "Test %d failed:\ninput:\n%s\nexpected:%s\n got:%s\n", i+1, tc.Input, tc.Output, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
