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

func simulate(s string, firstBias bool) int {
	first, second := 0, 0
	remainingFirst, remainingSecond := 5, 5
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			remainingFirst--
			if s[i] == '1' || (s[i] == '?' && firstBias) {
				first++
			}
		} else {
			remainingSecond--
			if s[i] == '1' || (s[i] == '?' && !firstBias) {
				second++
			}
		}
		if first > second+remainingSecond {
			return i + 1
		}
		if second > first+remainingFirst {
			return i + 1
		}
	}
	return 10
}

func solve(s string) int {
	res1 := simulate(s, true)
	res2 := simulate(s, false)
	if res1 < res2 {
		return res1
	}
	return res2
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var t int
	fmt.Fscan(rdr, &t)
	res := make([]string, t)
	for i := 0; i < t; i++ {
		var s string
		fmt.Fscan(rdr, &s)
		res[i] = fmt.Sprint(solve(s))
	}
	return strings.Join(res, "\n")
}

func randString() string {
	letters := []byte{'0', '1', '?'}
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateCases() []testCase {
	rand.Seed(3)
	cases := []testCase{}
	fixed := []string{"1111111111", "0000000000", "??????????", "1010101010"}
	for _, s := range fixed {
		inp := fmt.Sprintf("1\n%s\n", s)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		s := randString()
		inp := fmt.Sprintf("1\n%s\n", s)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierC.go <binary>")
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
