package main

import (
	"bufio"
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

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var t int
	fmt.Fscan(rdr, &t)
	res := make([]string, t)
	for i := 0; i < t; i++ {
		var n int64
		fmt.Fscan(rdr, &n)
		ans := n / 10
		if n%10 == 9 {
			ans++
		}
		res[i] = fmt.Sprint(ans)
	}
	return strings.Join(res, "\n")
}

func generateCases() []testCase {
	rand.Seed(1)
	cases := []testCase{}
	fixed := []int64{0, 1, 9, 10, 11, 18, 19, 20, 99, 100}
	for _, n := range fixed {
		inp := fmt.Sprintf("1\n%d\n", n)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Int63n(1_000_000_000) + 1
		inp := fmt.Sprintf("1\n%d\n", n)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
