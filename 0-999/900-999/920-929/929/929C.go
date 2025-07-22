package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Player struct {
	num  int
	role int // 0-goalie, 1-defender, 2-forward
}

func comb(n, k int) int64 {
	if n < k || k < 0 {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := 1; i <= k; i++ {
		res = res * int64(n-k+i) / int64(i)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var g, d, f int
	if _, err := fmt.Fscan(reader, &g, &d, &f); err != nil {
		return
	}

	players := make([]Player, g+d+f)
	for i := 0; i < g; i++ {
		fmt.Fscan(reader, &players[i].num)
		players[i].role = 0
	}
	for i := 0; i < d; i++ {
		fmt.Fscan(reader, &players[g+i].num)
		players[g+i].role = 1
	}
	for i := 0; i < f; i++ {
		fmt.Fscan(reader, &players[g+d+i].num)
		players[g+d+i].role = 2
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].num < players[j].num
	})

	n := len(players)
	prefixG := make([]int, n+1)
	prefixD := make([]int, n+1)
	prefixF := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefixG[i+1] = prefixG[i]
		prefixD[i+1] = prefixD[i]
		prefixF[i+1] = prefixF[i]
		switch players[i].role {
		case 0:
			prefixG[i+1]++
		case 1:
			prefixD[i+1]++
		case 2:
			prefixF[i+1]++
		}
	}

	var ans int64
	j := 0
	for i := 0; i < n; i++ {
		if j < i+1 {
			j = i + 1
		}
		for j < n && players[j].num <= 2*players[i].num {
			j++
		}
		numG := prefixG[j] - prefixG[i]
		numD := prefixD[j] - prefixD[i]
		numF := prefixF[j] - prefixF[i]
		switch players[i].role {
		case 0:
			if numD >= 2 && numF >= 3 {
				ans += comb(numD, 2) * comb(numF, 3)
			}
		case 1:
			if numG >= 1 && numD >= 2 && numF >= 3 {
				ans += int64(numG) * comb(numD-1, 1) * comb(numF, 3)
			}
		case 2:
			if numG >= 1 && numD >= 2 && numF >= 3 {
				ans += int64(numG) * comb(numD, 2) * comb(numF-1, 2)
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
