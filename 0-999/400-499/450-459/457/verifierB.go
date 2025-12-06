package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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
	rdr := strings.NewReader(input)
	var m, n int
	if _, err := fmt.Fscan(rdr, &m, &n); err != nil {
		return ""
	}
	a := make([]int, m)
	var sumA uint64
	for i := 0; i < m; i++ {
		fmt.Fscan(rdr, &a[i])
		sumA += uint64(a[i])
	}
	b := make([]int, n)
	var sumB uint64
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &b[i])
		sumB += uint64(b[i])
	}

	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	sort.Sort(sort.Reverse(sort.IntSlice(b)))

	ans := ^uint64(0)

	var currentSumA uint64
	for k := 1; k <= m; k++ {
		currentSumA += uint64(a[k-1])
		ops := (sumA - currentSumA) + uint64(k)*sumB
		if ops < ans {
			ans = ops
		}
	}

	var currentSumB uint64
	for l := 1; l <= n; l++ {
		currentSumB += uint64(b[l-1])
		ops := (sumB - currentSumB) + uint64(l)*sumA
		if ops < ans {
			ans = ops
		}
	}

	return fmt.Sprintf("%d", ans)
}

func generateCases() []testCase {
	rand.Seed(2)
	cases := []testCase{}
	fixed := []string{
		"1 1\n5\n7\n",
		"2 2\n1 2\n3 4\n",
		"3 1\n10 20 30\n5\n",
		"4 3\n67 38 62 93\n11 1 35\n",
	}
	for _, f := range fixed {
		cases = append(cases, testCase{f, compute(f)})
	}
	for len(cases) < 100 {
		m := rand.Intn(4) + 1
		n := rand.Intn(4) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", m, n)
		for i := 0; i < m; i++ {
			fmt.Fprintf(&sb, "%d", rand.Intn(100)+1)
			if i+1 < m {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d", rand.Intn(100)+1)
			if i+1 < n {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		inp := sb.String()
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
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:%s\nactual:%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}