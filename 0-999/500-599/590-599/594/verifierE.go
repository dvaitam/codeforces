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

func reverseString(str string) string {
	b := []byte(str)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func minString(a, b string) string {
	if a == "" {
		return b
	}
	if b == "" {
		return a
	}
	if a < b {
		return a
	}
	return b
}

func solveE(s string, k int) string {
	memo := make(map[[2]int]string)
	var dfs func(int, int) string
	dfs = func(pos, k int) string {
		if pos == len(s) {
			return ""
		}
		if k == 1 {
			rem := s[pos:]
			rev := reverseString(rem)
			if rem < rev {
				return rem
			}
			return rev
		}
		key := [2]int{pos, k}
		if v, ok := memo[key]; ok {
			return v
		}
		best := ""
		n := len(s)
		for i := pos + 1; i <= n; i++ {
			if k-1 > n-i {
				continue
			}
			part := s[pos:i]
			rest := dfs(i, k-1)
			best = minString(best, part+rest)
			best = minString(best, reverseString(part)+rest)
		}
		memo[key] = best
		return best
	}
	return dfs(0, k)
}

func run(binary, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for t := 0; t < 100; t++ {
		n := rand.Intn(6) + 1
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		s := string(b)
		k := rand.Intn(n) + 1
		input := fmt.Sprintf("%s\n%d\n", s, k)
		expected := solveE(s, k)
		got, err := run(bin, input)
		if err != nil {
			fmt.Println("test", t, "runtime error:", err)
			fmt.Println("output:\n" + got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Println("test", t, "failed")
			fmt.Println("input:\n" + input)
			fmt.Println("expected:\n" + expected)
			fmt.Println("got:\n" + got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
