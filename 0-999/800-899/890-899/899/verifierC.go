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
	if n%2 == 0 {
		if n == 2 {
			return "1\n1 1"
		}
		var sb strings.Builder
		if n%4 == 0 {
			sb.WriteString("0\n")
		} else {
			sb.WriteString("1\n")
		}
		size := n / 2
		fmt.Fprintf(&sb, "%d ", size)
		flag := true
		for i := 1; i <= n; i += 2 {
			if flag {
				fmt.Fprintf(&sb, "%d ", i)
			} else {
				fmt.Fprintf(&sb, "%d ", i+1)
			}
			flag = !flag
		}
		return strings.TrimSpace(sb.String())
	}
	if n == 3 {
		return "0\n1 3"
	}
	var sb strings.Builder
	if n%4 == 1 {
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d ", n/2)
	} else {
		sb.WriteString("0\n")
		fmt.Fprintf(&sb, "%d 1 ", n/2+1)
	}
	flag := true
	for i := 2; i <= n; i += 2 {
		if flag {
			fmt.Fprintf(&sb, "%d ", i)
		} else {
			fmt.Fprintf(&sb, "%d ", i+1)
		}
		flag = !flag
	}
	return strings.TrimSpace(sb.String())
}

func generateTests() []test {
	rand.Seed(8993)
	var tests []test
	fixed := []string{
		"4\n",
		"3\n",
		"5\n",
	}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rand.Intn(60) + 2
		inp := fmt.Sprintf("%d\n", n)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
