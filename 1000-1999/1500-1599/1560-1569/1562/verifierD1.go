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

func solveCase(n, q int, s string, queries [][2]int) string {
	p := make([]int, n+1)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			if s[i] == '+' {
				p[i+1] = p[i] + 1
			} else {
				p[i+1] = p[i] - 1
			}
		} else {
			if s[i] == '+' {
				p[i+1] = p[i] - 1
			} else {
				p[i+1] = p[i] + 1
			}
		}
	}
	var out strings.Builder
	for _, qr := range queries {
		l, r := qr[0], qr[1]
		c := p[r] - p[l-1]
		if c == 0 {
			out.WriteString("0\n")
		} else if c%2 != 0 {
			out.WriteString("1\n")
		} else {
			out.WriteString("2\n")
		}
	}
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	str := make([]byte, n)
	for i := range str {
		if rng.Intn(2) == 0 {
			str[i] = '+'
		} else {
			str[i] = '-'
		}
	}
	s := string(str)
	sb.WriteString(s)
	sb.WriteByte('\n')
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		queries[i] = [2]int{l, r}
		sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
	}
	input := sb.String()
	expected := solveCase(n, q, s, queries)
	return input, expected
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
