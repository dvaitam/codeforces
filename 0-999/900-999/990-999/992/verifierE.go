package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type query struct {
	p int
	x int
}

type testCaseE struct {
	input    string
	expected []int
}

func computeE(a []int, qs []query) []int {
	n := len(a)
	res := make([]int, len(qs))
	for idx, q := range qs {
		a[q.p-1] = q.x
		sum := 0
		ans := -1
		for i := 0; i < n; i++ {
			if a[i] == sum {
				ans = i + 1
				break
			}
			sum += a[i]
		}
		res[idx] = ans
	}
	return res
}

func generateCaseE(rng *rand.Rand) testCaseE {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10)
	}
	queries := make([]query, q)
	for i := 0; i < q; i++ {
		p := rng.Intn(n) + 1
		x := rng.Intn(10)
		queries[i] = query{p: p, x: x}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu.p, qu.x))
	}
	return testCaseE{input: sb.String(), expected: computeE(a, queries)}
}

func parseOutputE(out string, q int) ([]int, error) {
	lines := strings.Fields(out)
	if len(lines) != q {
		return nil, fmt.Errorf("expected %d lines", q)
	}
	res := make([]int, q)
	for i, s := range lines {
		var v int
		if _, err := fmt.Sscan(s, &v); err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		tc := generateCaseE(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, tc.input)
			os.Exit(1)
		}
		vals, err := parseOutputE(out, len(tc.expected))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\ninput:\n%s", i, tc.input)
			os.Exit(1)
		}
		for j, v := range vals {
			if v != tc.expected[j] {
				fmt.Fprintf(os.Stderr, "case %d: expected %v got %v\ninput:\n%s", i, tc.expected, vals, tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
