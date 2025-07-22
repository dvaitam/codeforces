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
	input string
	ans   int
}

func expected(a, b []int) int {
	n := len(a)
	pos := make([]int, n+1)
	for i, v := range a {
		pos[v] = i
	}
	length := 1
	cur := pos[b[n-1]]
	for i := n - 2; i >= 0; i-- {
		p := pos[b[i]]
		if p < cur {
			length++
			cur = p
		} else {
			break
		}
	}
	return n - length
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	a := rng.Perm(n)
	b := rng.Perm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v + 1))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v + 1))
	}
	sb.WriteByte('\n')
	aa := make([]int, n)
	bb := make([]int, n)
	for i := 0; i < n; i++ {
		aa[i] = a[i] + 1
		bb[i] = b[i] + 1
	}
	return testCase{input: sb.String(), ans: expected(aa, bb)}
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
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != tc.ans {
		return fmt.Errorf("expected %d got %d", tc.ans, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{input: "1\n1\n1\n", ans: 0},
		{input: "3\n1 2 3\n1 2 3\n", ans: 0},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
