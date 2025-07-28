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

func canForm(s, t string) bool {
	n := len(s)
	m := len(t)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			k := j - i + 1
			if k > m {
				break
			}
			ok := true
			for p := 0; p < k; p++ {
				if s[i+p] != t[p] {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
			rem := m - k
			if rem == 0 {
				return true
			}
			if j-rem < 0 {
				continue
			}
			ok2 := true
			for p := 0; p < rem; p++ {
				if s[j-1-p] != t[k+p] {
					ok2 = false
					break
				}
			}
			if ok2 {
				return true
			}
		}
	}
	return false
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var q int
	fmt.Fscan(rdr, &q)
	res := make([]string, q)
	for i := 0; i < q; i++ {
		var s, t string
		fmt.Fscan(rdr, &s)
		fmt.Fscan(rdr, &t)
		if canForm(s, t) {
			res[i] = "YES"
		} else {
			res[i] = "NO"
		}
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
	rand.Seed(2)
	cases := []testCase{}
	fixed := []struct{ s, t string }{
		{"a", "a"},
		{"ab", "ba"},
		{"abc", "abc"},
		{"abc", "cba"},
		{"abcdef", "cdedcb"},
	}
	for _, f := range fixed {
		inp := fmt.Sprintf("1\n%s\n%s\n", f.s, f.t)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Intn(6) + 1
		m := rand.Intn(2*n-1) + 1
		s := randString(n)
		t := randString(m)
		inp := fmt.Sprintf("1\n%s\n%s\n", s, t)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB.go <binary>")
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
