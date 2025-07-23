package main

import (
	"bufio"
	"fmt"
	"os"
)

func best(p, cards, k int, h []int) int {
	limit := p * k
	if cards > limit {
		cards = limit
	}
	dp := make([]int, cards+1)
	for i := 0; i < p; i++ {
		next := make([]int, cards+1)
		for j := 0; j <= cards; j++ {
			for t := 0; t <= k && t <= j; t++ {
				val := dp[j-t] + h[t]
				if val > next[j] {
					next[j] = val
				}
			}
		}
		dp = next
	}
	res := 0
	for j := 0; j <= cards; j++ {
		if dp[j] > res {
			res = dp[j]
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	total := n * k
	cards := make([]int, total)
	for i := 0; i < total; i++ {
		fmt.Fscan(reader, &cards[i])
	}
	fav := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &fav[i])
	}
	h := make([]int, k+1)
	h[0] = 0
	for i := 1; i <= k; i++ {
		fmt.Fscan(reader, &h[i])
	}

	type pair struct {
		players int
		count   int
	}

	mp := make(map[int]*pair)
	for _, f := range fav {
		if mp[f] == nil {
			mp[f] = &pair{}
		}
		mp[f].players++
	}
	for _, c := range cards {
		if mp[c] == nil {
			mp[c] = &pair{}
		}
		mp[c].count++
	}

	ans := 0
	for _, p := range mp {
		if p.players > 0 {
			ans += best(p.players, p.count, k, h)
		}
	}
	fmt.Fprintln(writer, ans)
}
