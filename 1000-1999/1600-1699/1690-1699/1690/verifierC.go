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

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func solveCase(s, f []int) []int {
	n := len(s)
	ans := make([]int, n)
	prev := 0
	for i := 0; i < n; i++ {
		start := s[i]
		if prev > start {
			start = prev
		}
		ans[i] = f[i] - start
		prev = f[i]
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(44))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(10) + 1
		s := make([]int, n)
		f := make([]int, n)
		cur := rng.Intn(5)
		for i := 0; i < n; i++ {
			cur += rng.Intn(5) + 1
			s[i] = cur
		}
		cur = s[0] + rng.Intn(5) + 1
		for i := 0; i < n; i++ {
			if i > 0 {
				if cur < s[i] {
					cur = s[i]
				}
				cur += rng.Intn(5) + 1
			}
			f[i] = cur
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", s[i])
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", f[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := solveCase(s, f)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) < n {
			fmt.Fprintf(os.Stderr, "case %d wrong output size: got %q\n", tc+1, out)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			var v int
			fmt.Sscanf(fields[i], "%d", &v)
			if v != exp[i] {
				fmt.Fprintf(os.Stderr, "case %d failed at index %d: expected %d got %d\ninput:\n%s", tc+1, i, exp[i], v, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
