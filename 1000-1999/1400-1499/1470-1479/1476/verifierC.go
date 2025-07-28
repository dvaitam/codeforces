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

type testCase struct {
	input  string
	output string
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveCase(c, a, b []int) int {
	n := len(c)
	cur, ans := 0, 0
	for i := 1; i < n; i++ {
		diff := abs(a[i] - b[i])
		if diff == 0 {
			cur = 0
		} else {
			cur = max(diff, cur+c[i-1]-diff)
		}
		ans = max(ans, cur+c[i])
	}
	return ans + 1
}

func buildCase(c, a, b []int) testCase {
	n := len(c)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range c {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	out := fmt.Sprintf("%d\n", solveCase(c, a, b))
	return testCase{input: sb.String(), output: out}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 2
	c := make([]int, n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		c[i] = rng.Intn(10) + 2
		a[i] = rng.Intn(c[i]-1) + 1
		b[i] = rng.Intn(c[i]-1) + 1
	}
	return buildCase(c, a, b)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.output)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	cases = append(cases, buildCase([]int{2, 2}, []int{1, 1}, []int{1, 2}))
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
