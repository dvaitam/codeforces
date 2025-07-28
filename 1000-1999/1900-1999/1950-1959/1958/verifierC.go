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

func minSplits(n int64, k int64) int64 {
	var steps int64
	for k > 0 && k < (int64(1)<<n) {
		half := int64(1) << (n - 1)
		if k > half {
			k -= half
		} else {
			k = half - k
		}
		steps++
		n--
	}
	return steps
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(44)
	const t = 100
	ns := make([]int64, t)
	ks := make([]int64, t)
	for i := 0; i < t; i++ {
		ns[i] = int64(rand.Intn(60) + 1)
		maxK := int64(1)<<ns[i] - 1
		ks[i] = rand.Int63n(maxK) + 1
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

	outs := strings.Fields(string(outBytes))
	if len(outs) != t {
		fmt.Printf("expected %d lines, got %d\n", t, len(outs))
		os.Exit(1)
	}
	for i, s := range outs {
		var got int64
		fmt.Sscan(s, &got)
		want := minSplits(ns[i], ks[i])
		if got != want {
			fmt.Printf("mismatch on test %d: n=%d k=%d expected %d got %d\n", i+1, ns[i], ks[i], want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
