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

func expected(k, m int64) int64 {
	mod := m % (3 * k)
	if mod >= 2*k {
		return 0
	}
	return 2*k - mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(43)
	const t = 100
	ks := make([]int64, t)
	ms := make([]int64, t)
	for i := 0; i < t; i++ {
		ks[i] = rand.Int63n(1e8) + 1
		ms[i] = rand.Int63n(1e9) + 1
	}

	var input strings.Builder
	fmt.Fprintln(&input, t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d %d\n", ks[i], ms[i])
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
		var got int64
		fmt.Sscan(s, &got)
		want := expected(ks[i], ms[i])
		if got != want {
			fmt.Printf("mismatch on test %d: k=%d m=%d expected %d got %d\n", i+1, ks[i], ms[i], want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
