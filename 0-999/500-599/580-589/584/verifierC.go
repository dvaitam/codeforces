package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func note(x, y byte) byte {
	for z := byte('a'); z <= 'z'; z++ {
		if z != x && z != y {
			return z
		}
	}
	return 'a'
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(s1, s2 string, t int) (string, bool) {
	n := len(s1)
	diffIdx := make([]int, 0, n)
	sameIdx := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if s1[i] != s2[i] {
			diffIdx = append(diffIdx, i)
		} else {
			sameIdx = append(sameIdx, i)
		}
	}

	d := len(diffIdx)
	s := n - d

	low := max(0, d-t)
	high := min(d/2, n-t)
	if low > high {
		return "", false
	}
	x := low
	y := t - d + x

	ans := make([]byte, n)
	for i := 0; i < x; i++ {
		idx := diffIdx[i]
		ans[idx] = s1[idx]
	}
	for i := x; i < 2*x; i++ {
		idx := diffIdx[i]
		ans[idx] = s2[idx]
	}
	for i := 2 * x; i < d; i++ {
		idx := diffIdx[i]
		ans[idx] = note(s1[idx], s2[idx])
	}
	for i := 0; i < y; i++ {
		idx := sameIdx[i]
		ans[idx] = note(s1[idx], 0)
	}
	for i := y; i < s; i++ {
		idx := sameIdx[i]
		ans[idx] = s1[idx]
	}
	return string(ans), true
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func diff(a, b string) int {
	cnt := 0
	for i := range a {
		if a[i] != b[i] {
			cnt++
		}
	}
	return cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type tc struct {
		n      int
		t      int
		s1, s2 string
	}
	cases := make([]tc, 0, 100)
	// some edge cases
	cases = append(cases, tc{1, 0, "a", "a"})
	cases = append(cases, tc{1, 1, "a", "b"})
	cases = append(cases, tc{2, 1, "aa", "aa"})

	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		t := rng.Intn(n + 1)
		var b1, b2 strings.Builder
		for i := 0; i < n; i++ {
			b1.WriteByte(byte('a' + rng.Intn(26)))
			b2.WriteByte(byte('a' + rng.Intn(26)))
		}
		cases = append(cases, tc{n, t, b1.String(), b2.String()})
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n%s\n%s\n", tc.n, tc.t, tc.s1, tc.s2)
		_, ok := solve(tc.s1, tc.s2, tc.t)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		if !ok {
			if out != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:%s", idx+1, out, input)
				os.Exit(1)
			}
			continue
		}
		if out == "-1" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected string got -1\ninput:%s", idx+1, input)
			os.Exit(1)
		}
		if len(out) != tc.n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected length %d got %d\n", idx+1, tc.n, len(out))
			os.Exit(1)
		}
		if diff(tc.s1, out) != tc.t || diff(tc.s2, out) != tc.t {
			fmt.Fprintf(os.Stderr, "case %d failed: output doesn't differ by t\ninput:%soutput:%s\n", idx+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
