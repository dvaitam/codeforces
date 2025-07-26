package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func isNice(a []int, k int) bool {
	n := len(a)
	next := make([]int, n)
	stack := make([]int, 0, n)
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			next[i] = n
		} else {
			next[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	for i := 0; i < n; i++ {
		if next[i]-i > k {
			return false
		}
	}
	return true
}

func solveCase(a []int, k int) string {
	n := len(a)
	if k >= n {
		return "YES"
	}
	if isNice(a, k) {
		return "YES"
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if a[i] <= a[j] {
				continue
			}
			a[i], a[j] = a[j], a[i]
			if isNice(a, k) {
				a[i], a[j] = a[j], a[i]
				return "YES"
			}
			a[i], a[j] = a[j], a[i]
		}
	}
	return "NO"
}

func runCase(bin string, a []int, k int) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(a), k)
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(a, k)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(n) + 1
		a := make([]int, n)
		for j := range a {
			a[j] = rand.Intn(100) + 1
		}
		if err := runCase(bin, a, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
