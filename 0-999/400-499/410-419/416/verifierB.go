package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var m, n int
	if _, err := fmt.Fscan(in, &m, &n); err != nil {
		return ""
	}
	dp := make([]int, n)
	res := make([]int, m)
	for i := 0; i < m; i++ {
		prev := 0
		for j := 0; j < n; j++ {
			var t int
			fmt.Fscan(in, &t)
			start := dp[j]
			if prev > start {
				start = prev
			}
			finish := start + t
			dp[j] = finish
			prev = finish
		}
		res[i] = prev
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	m := rng.Intn(5) + 1
	n := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", m, n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			t := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("%d", t))
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	exp := solve(input)
	return input, strings.TrimSpace(exp)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
