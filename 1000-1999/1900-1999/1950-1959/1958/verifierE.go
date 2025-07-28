package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func construct(n, k int) []int {
	if k == 1 {
		res := make([]int, n)
		for i := 0; i < n; i++ {
			res[i] = i + 1
		}
		return res
	}
	if k == 2 {
		if n < 3 {
			return nil
		}
		res := []int{n - 1, 1, n}
		used := map[int]bool{n - 1: true, 1: true, n: true}
		for i := 2; i <= n-2; i++ {
			if !used[i] {
				res = append(res, i)
			}
		}
		return res
	}
	if k == 3 {
		if n < 5 {
			return nil
		}
		res := []int{n - 1, 1, n - 2, 2, n}
		used := map[int]bool{n - 1: true, 1: true, n - 2: true, 2: true, n: true}
		for i := 3; i <= n; i++ {
			if !used[i] {
				res = append(res, i)
			}
		}
		return res
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(46)
	const t = 100
	ns := make([]int, t)
	ks := make([]int, t)
	for i := 0; i < t; i++ {
		ns[i] = rand.Intn(10) + 2
		ks[i] = rand.Intn(ns[i]-1) + 1
	}

	var input strings.Builder
	fmt.Fprintln(&input, t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d %d\n", ns[i], ks[i])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input.String())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}

	outLines := strings.Split(strings.TrimSpace(string(outBytes)), "\n")
	if len(outLines) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outLines))
		os.Exit(1)
	}

	for i, line := range outLines {
		want := construct(ns[i], ks[i])
		if want == nil {
			if strings.TrimSpace(line) != "-1" {
				fmt.Printf("test %d expected -1 got %s\n", i+1, line)
				os.Exit(1)
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != ns[i] {
			fmt.Printf("test %d expected %d numbers got %d\n", i+1, ns[i], len(parts))
			os.Exit(1)
		}
		for j, s := range parts {
			var got int
			fmt.Sscan(s, &got)
			if got != want[j] {
				fmt.Printf("test %d mismatch at position %d expected %d got %d\n", i+1, j+1, want[j], got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}
