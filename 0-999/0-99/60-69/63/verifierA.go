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

func runBinary(bin string, input string) (string, error) {
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

func generateCases() []testCase {
	rand.Seed(1)
	cases := make([]testCase, 100)
	statuses := []string{"rat", "woman", "child", "man"}
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		buf := bytes.Buffer{}
		fmt.Fprintln(&buf, n)
		names := make([]string, n)
		stat := make([]string, n)
		capPos := rand.Intn(n)
		for j := 0; j < n; j++ {
			names[j] = fmt.Sprintf("N%02d_%02d", i, j)
			if j == capPos {
				stat[j] = "captain"
			} else {
				stat[j] = statuses[rand.Intn(len(statuses))]
			}
			fmt.Fprintf(&buf, "%s %s\n", names[j], stat[j])
		}
		// expected output
		var exp bytes.Buffer
		for idx := 0; idx < n; idx++ {
			if stat[idx] == "rat" {
				fmt.Fprintln(&exp, names[idx])
			}
		}
		for idx := 0; idx < n; idx++ {
			if stat[idx] == "woman" || stat[idx] == "child" {
				fmt.Fprintln(&exp, names[idx])
			}
		}
		for idx := 0; idx < n; idx++ {
			if stat[idx] == "man" {
				fmt.Fprintln(&exp, names[idx])
			}
		}
		for idx := 0; idx < n; idx++ {
			if stat[idx] == "captain" {
				fmt.Fprintln(&exp, names[idx])
				break
			}
		}
		cases[i] = testCase{input: buf.String(), expected: strings.TrimSpace(exp.String())}
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
