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

func solve(n int, a, b []int64, queries [][2]int) []int64 {
	prefixB := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefixB[i] = prefixB[i-1] + b[i]
	}
	ans := make([]int64, len(queries))
	for qi, q := range queries {
		l, r := q[0], q[1]
		var moves int64
		if l < r {
			for j := l + 1; j <= r; j++ {
				power := prefixB[j-1] - prefixB[l-1]
				moves += (a[j] + power - 1) / power
			}
		}
		ans[qi] = moves
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(51)
	n := rand.Intn(10) + 1
	qn := 100
	a := make([]int64, n+1)
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = int64(rand.Intn(20) + 1)
	}
	for i := 1; i <= n; i++ {
		b[i] = int64(rand.Intn(20) + 1)
	}
	queries := make([][2]int, qn)
	for i := 0; i < qn; i++ {
		l := rand.Intn(n) + 1
		r := rand.Intn(n) + 1
		if l > r {
			l, r = r, l
		}
		queries[i] = [2]int{l, r}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, qn)
	for i := 1; i <= n; i++ {
		if i > 1 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, a[i])
	}
	input.WriteByte('\n')
	for i := 1; i <= n; i++ {
		if i > 1 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, b[i])
	}
	input.WriteByte('\n')
	for _, q := range queries {
		fmt.Fprintf(&input, "%d %d\n", q[0], q[1])
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

	outs := strings.Fields(strings.TrimSpace(string(outBytes)))
	if len(outs) != qn {
		fmt.Printf("expected %d lines, got %d\n", qn, len(outs))
		os.Exit(1)
	}
	want := solve(n, a, b, queries)
	for i, s := range outs {
		var got int64
		fmt.Sscan(s, &got)
		if got != want[i] {
			fmt.Printf("mismatch on query %d expected %d got %d\n", i+1, want[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
