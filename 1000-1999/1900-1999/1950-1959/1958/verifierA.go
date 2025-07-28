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

// reference implementation of problem A
func minOnes(n int) int {
	for ones := 0; ones <= n; ones++ {
		m := n - ones
		if m < 0 {
			break
		}
		for a := 0; 3*a <= m; a++ {
			if (m-3*a)%5 == 0 {
				return ones
			}
		}
	}
	return n
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(42)
	const t = 100
	ns := make([]int, t)
	for i := range ns {
		ns[i] = rand.Intn(100) + 1 // 1..100
	}

	var input strings.Builder
	fmt.Fprintln(&input, t)
	for _, n := range ns {
		fmt.Fprintln(&input, n)
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
	outs := strings.Fields(string(outBytes))
	if len(outs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outs))
		os.Exit(1)
	}

	for i, s := range outs {
		var got int
		fmt.Sscan(s, &got)
		want := minOnes(ns[i])
		if got != want {
			fmt.Printf("mismatch on test %d: n=%d expected %d got %d\n", i+1, ns[i], want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
