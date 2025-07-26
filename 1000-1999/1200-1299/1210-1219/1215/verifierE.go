package main

import (
	"bytes"
	"context"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(arr []int) int64 {
	n := len(arr)
	freq := make([]int, 20)
	for _, v := range arr {
		freq[v-1]++
	}
	idx := make([]int, 20)
	m := 0
	for c := 0; c < 20; c++ {
		if freq[c] > 0 {
			idx[c] = m
			m++
		} else {
			idx[c] = -1
		}
	}
	comp := make([]int, n)
	for i := 0; i < n; i++ {
		comp[i] = idx[arr[i]-1]
	}
	cross := make([][]int64, m)
	for i := range cross {
		cross[i] = make([]int64, m)
	}
	cnt := make([]int64, m)
	for _, c := range comp {
		for j := 0; j < m; j++ {
			cross[j][c] += cnt[j]
		}
		cnt[c]++
	}
	size := 1 << m
	cost := make([][]int64, m)
	for i := 0; i < m; i++ {
		cost[i] = make([]int64, size)
		for mask := 1; mask < size; mask++ {
			b := mask & -mask
			j := bits.TrailingZeros(uint(b))
			prev := mask ^ b
			cost[i][mask] = cost[i][prev] + cross[i][j]
		}
	}
	const INF int64 = 1 << 60
	dp := make([]int64, size)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for mask := 0; mask < size; mask++ {
		cur := dp[mask]
		if cur == INF {
			continue
		}
		for c := 0; c < m; c++ {
			if mask&(1<<c) == 0 {
				nm := mask | (1 << c)
				val := cur + cost[c][mask]
				if val < dp[nm] {
					dp[nm] = val
				}
			}
		}
	}
	return dp[size-1]
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
	binary := os.Args[1]

	rand.Seed(46)
	const T = 100
	for i := 0; i < T; i++ {
		n := rand.Intn(7) + 3
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rand.Intn(6) + 1
		}
		input := fmt.Sprintf("%d\n", n)
		for j, v := range arr {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		expect := solve(arr)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\ninput:%s\n", i+1, err, input)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscanf(out, "%d", &got); err != nil {
			fmt.Printf("test %d: failed to parse output '%s'\n", i+1, out)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\ninput:%s\n", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
