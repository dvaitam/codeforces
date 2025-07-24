package main

import (
	"bufio"
	"fmt"
	"os"
)

func shiftOr(dp []uint64, shift int) {
	ws := shift / 64
	bs := shift % 64
	if ws >= len(dp) {
		return
	}
	if bs == 0 {
		for i := len(dp) - 1; i >= ws; i-- {
			dp[i] |= dp[i-ws]
		}
	} else {
		for i := len(dp) - 1; i > ws; i-- {
			dp[i] |= dp[i-ws]<<bs | dp[i-ws-1]>>(64-bs)
		}
		dp[ws] |= dp[0] << bs
	}
}

func dpSubset(k int, freq map[int]int) []uint64 {
	words := k/64 + 1
	dp := make([]uint64, words)
	dp[0] = 1
	for l, c := range freq {
		step := 1
		for c > 0 {
			if step > c {
				step = c
			}
			shift := l * step
			if shift <= k {
				shiftOr(dp, shift)
			}
			c -= step
			step <<= 1
		}
	}
	return dp
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}

	visited := make([]bool, n)
	var cycles []int
	for i := 0; i < n; i++ {
		if !visited[i] {
			x := i
			l := 0
			for !visited[x] {
				visited[x] = true
				x = p[x] - 1
				l++
			}
			cycles = append(cycles, l)
		}
	}

	freq := make(map[int]int)
	limit := 0
	odd := 0
	for _, l := range cycles {
		freq[l]++
		limit += l / 2
		if l%2 == 1 {
			odd++
		}
	}

	var maxFail int
	if k <= limit {
		maxFail = 2 * k
	} else if k <= limit+odd {
		maxFail = limit + k
	} else {
		maxFail = n
	}

	dp := dpSubset(k, freq)
	bit := (dp[k/64] >> (k % 64)) & 1
	minFail := k
	if bit == 0 {
		if k < n {
			minFail = k + 1
		} else {
			minFail = n
		}
	}

	fmt.Fprintf(writer, "%d %d\n", minFail, maxFail)
}
