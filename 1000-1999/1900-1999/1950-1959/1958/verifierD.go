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

func costBlock(arr []int64) int64 {
	m := len(arr)
	if m%2 == 0 {
		var cost int64
		for i := 0; i < m; i += 2 {
			cost += 2 * (arr[i] + arr[i+1])
		}
		return cost
	}
	pairOdd := make([]int64, m+1)
	pairEven := make([]int64, m+1)
	for i := 1; i < m; i++ {
		pairOdd[i] = pairOdd[i-1]
		pairEven[i] = pairEven[i-1]
		if i%2 == 1 {
			pairOdd[i] += 2 * (arr[i-1] + arr[i])
		} else {
			pairEven[i] += 2 * (arr[i-1] + arr[i])
		}
	}
	pairOdd[m] = pairOdd[m-1]
	pairEven[m] = pairEven[m-1]
	ans := int64(1<<62 - 1)
	for j := 1; j <= m; j += 2 {
		cost := pairOdd[j-1] + arr[j-1] + (pairEven[m] - pairEven[j])
		if cost < ans {
			ans = cost
		}
	}
	return ans
}

func expected(a []int64) int64 {
	n := len(a)
	var total int64
	for i := 0; i < n; {
		if a[i] == 0 {
			i++
			continue
		}
		start := i
		for i < n && a[i] > 0 {
			i++
		}
		total += costBlock(a[start:i])
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(45)
	const t = 100
	ns := make([]int, t)
	arrays := make([][]int64, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1 // keep small
		ns[i] = n
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			if rand.Intn(3) == 0 {
				arr[j] = 0
			} else {
				arr[j] = rand.Int63n(100) + 1
			}
		}
		arrays[i] = arr
	}

	var input strings.Builder
	fmt.Fprintln(&input, t)
	for i := 0; i < t; i++ {
		fmt.Fprintln(&input, ns[i])
		for j, v := range arrays[i] {
			if j > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, v)
		}
		input.WriteByte('\n')
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
		want := expected(arrays[i])
		if got != want {
			fmt.Printf("mismatch on test %d expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
