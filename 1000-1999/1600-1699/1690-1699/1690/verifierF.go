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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func minimalShift(chars []byte) int {
	m := len(chars)
	for d := 1; d <= m; d++ {
		if m%d != 0 {
			continue
		}
		ok := true
		for i := 0; i < m; i++ {
			if chars[i] != chars[(i+d)%m] {
				ok = false
				break
			}
		}
		if ok {
			return d
		}
	}
	return m
}

func solveCase(n int, s string, p []int) int64 {
	visited := make([]bool, n)
	ans := int64(1)
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		cycle := []int{}
		j := i
		for !visited[j] {
			visited[j] = true
			cycle = append(cycle, j)
			j = p[j]
		}
		chars := make([]byte, len(cycle))
		for idx, v := range cycle {
			chars[idx] = s[v]
		}
		d := minimalShift(chars)
		ans = lcm(ans, int64(d))
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(47))
	letters := []byte("abc")
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(8) + 1
		b := make([]byte, n)
		for i := range b {
			b[i] = letters[rng.Intn(len(letters))]
		}
		perm := rng.Perm(n)
		p := make([]int, n)
		for i, v := range perm {
			p[i] = v
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d\n%s\n", n, string(b))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", p[i]+1)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveCase(n, string(b), p)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", tc+1, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscanf(strings.TrimSpace(out), "%d", &got)
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", tc+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
