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

func solveCase(s, t string) string {
	i := len(s) - 1
	j := len(t) - 1
	for i >= 0 && j >= 0 {
		if s[i] == t[j] {
			i--
			j--
		} else {
			i -= 2
		}
	}
	if j < 0 {
		return "YES"
	}
	return "NO"
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var q int
	fmt.Fscan(rdr, &q)
	res := make([]string, q)
	for idx := 0; idx < q; idx++ {
		var s, t string
		fmt.Fscan(rdr, &s)
		fmt.Fscan(rdr, &t)
		res[idx] = solveCase(s, t)
	}
	return strings.Join(res, "\n")
}

func randString(n int) string {
	letters := []rune("abcde")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateCases() []testCase {
	rand.Seed(4)
	cases := []testCase{}
	fixed := []struct{ s, t string }{
		{"ababa", "ba"},
		{"ababa", "bb"},
		{"aaa", "aaaa"},
		{"aababa", "ababa"},
	}
	for _, f := range fixed {
		inp := fmt.Sprintf("1\n%s\n%s\n", f.s, f.t)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Intn(8) + 1
		m := rand.Intn(n) + 1
		s := randString(n)
		t := randString(m)
		inp := fmt.Sprintf("1\n%s\n%s\n", s, t)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD.go <binary>")
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
