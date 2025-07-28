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

func checkShift(p []int, n, m, k int) bool {
	visited := make([]bool, n)
	swaps := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			j := i
			cycleLen := 0
			for !visited[j] {
				visited[j] = true
				j = (p[j] + k) % n
				cycleLen++
			}
			swaps += cycleLen - 1
			if swaps > m {
				return false
			}
		}
	}
	return swaps <= m
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var t int
	fmt.Fscan(rdr, &t)
	outs := make([]string, t)
	for caseID := 0; caseID < t; caseID++ {
		var n, m int
		fmt.Fscan(rdr, &n, &m)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(rdr, &p[i])
			p[i]--
		}
		cnt := make([]int, n)
		for i := 0; i < n; i++ {
			k := (i - p[i] + n) % n
			cnt[k]++
		}
		var cand []int
		for k := 0; k < n; k++ {
			if n-cnt[k] <= 2*m {
				if checkShift(p, n, m, k) {
					cand = append(cand, k)
				}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(len(cand)))
		for _, k := range cand {
			sb.WriteByte(' ')
			sb.WriteString(fmt.Sprint(k))
		}
		outs[caseID] = sb.String()
	}
	return strings.Join(outs, "\n")
}

func randPerm(n int) []int {
	p := rand.Perm(n)
	for i := range p {
		p[i]++
	}
	return p
}

func generateCases() []testCase {
	rand.Seed(5)
	cases := []testCase{}
	fixed := []struct {
		n, m int
		p    []int
	}{
		{4, 1, []int{2, 3, 1, 4}},
		{3, 0, []int{1, 2, 3}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", f.n, f.m)
		for i, x := range f.p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(x))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Intn(6) + 3
		m := rand.Intn(n/3 + 1)
		p := randPerm(n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d\n", n, m)
		for i, x := range p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(x))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE.go <binary>")
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
