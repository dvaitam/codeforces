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
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateCases() []testCase {
	rand.Seed(1)
	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		r := rand.Intn(5) + 1
		c := rand.Intn(5) + 1
		var buf bytes.Buffer
		for row := 1; row <= 5; row++ {
			for col := 1; col <= 5; col++ {
				val := 0
				if row == r && col == c {
					val = 1
				}
				fmt.Fprintf(&buf, "%d", val)
				if col < 5 {
					buf.WriteByte(' ')
				}
			}
			buf.WriteByte('\n')
		}
		exp := fmt.Sprintf("%d", abs(r-3)+abs(c-3))
		cases[i] = testCase{input: buf.String(), expected: exp}
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA.go <binary>")
		os.Exit(1)
	}
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
