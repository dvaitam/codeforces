package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var n int
	fmt.Fscan(reader, &n)
	ones, twos := 0, 0
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		if v == 1 {
			ones++
		} else {
			twos++
		}
	}
	teams := min(ones, twos)
	ones -= teams
	teams += ones / 3
	return fmt.Sprintf("%d", teams)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func generateTests() []test {
	rand.Seed(8991)
	var tests []test
	fixed := []string{
		"4\n1 1 2 2\n",
		"7\n2 1 2 1 1 2 1\n",
	}
	for _, in := range fixed {
		tests = append(tests, test{in, solve(in)})
	}
	for len(tests) < 100 {
		n := rand.Intn(100) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			if rand.Intn(2) == 0 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('2')
			}
		}
		sb.WriteByte('\n')
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
