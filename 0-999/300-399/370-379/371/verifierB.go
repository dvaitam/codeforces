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

func factor(x int64) (e2, e3, e5 int, rem int64) {
	rem = x
	for rem%2 == 0 {
		rem /= 2
		e2++
	}
	for rem%3 == 0 {
		rem /= 3
		e3++
	}
	for rem%5 == 0 {
		rem /= 5
		e5++
	}
	return
}

func solve(input string) string {
	reader := strings.NewReader(input)
	var a, b int64
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return ""
	}
	a2, a3, a5, ra := factor(a)
	b2, b3, b5, rb := factor(b)
	if ra != rb {
		return "-1"
	}
	ops := abs(a2-b2) + abs(a3-b3) + abs(a5-b5)
	return fmt.Sprintf("%d", ops)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateTests() []test {
	rand.Seed(43)
	var tests []test
	fixed := []struct{ a, b int64 }{
		{1, 1}, {10, 10}, {10, 1}, {100, 250}, {45, 75},
	}
	for _, f := range fixed {
		inp := fmt.Sprintf("%d %d\n", f.a, f.b)
		tests = append(tests, test{inp, solve(inp)})
	}
	for len(tests) < 100 {
		a := rand.Int63n(1000000) + 1
		b := rand.Int63n(1000000) + 1
		inp := fmt.Sprintf("%d %d\n", a, b)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sexpected:%s\n got:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
