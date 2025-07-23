package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Group struct {
	length int
	costs  []int
}

func sumSmallest(freq []int, k int) int {
	if k <= 0 {
		return 0
	}
	sum := 0
	for cost := 1; cost < len(freq) && k > 0; cost++ {
		if freq[cost] > 0 {
			take := freq[cost]
			if take > k {
				take = k
			}
			sum += take * cost
			k -= take
		}
	}
	return sum
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	legs := make([]struct{ l, d int }, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &legs[i].l)
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &legs[i].d)
	}
	sort.Slice(legs, func(i, j int) bool { return legs[i].l < legs[j].l })

	groups := make([]Group, 0)
	for i := 0; i < n; {
		j := i
		for j < n && legs[j].l == legs[i].l {
			j++
		}
		g := Group{length: legs[i].l, costs: make([]int, j-i)}
		for k := i; k < j; k++ {
			g.costs[k-i] = legs[k].d
		}
		groups = append(groups, g)
		i = j
	}

	m := len(groups)
	suffixCost := make([]int, m+1)
	for i := m - 1; i >= 0; i-- {
		sum := 0
		for _, v := range groups[i].costs {
			sum += v
		}
		suffixCost[i] = suffixCost[i+1] + sum
	}

	freq := make([]int, 201)
	prefixCost := 0
	countShorter := 0
	ans := int(1<<31 - 1)

	for idx, g := range groups {
		sort.Ints(g.costs)
		c := len(g.costs)
		prefixGroup := make([]int, c+1)
		for i := 0; i < c; i++ {
			prefixGroup[i+1] = prefixGroup[i] + g.costs[i]
		}
		groupTotal := prefixGroup[c]
		costLonger := suffixCost[idx+1]
		for x := 1; x <= c; x++ {
			keepOthers := x - 1
			if keepOthers > countShorter {
				keepOthers = countShorter
			}
			costKeepOthers := sumSmallest(freq, keepOthers)
			costRemoveShorter := prefixCost - costKeepOthers
			costRemoveGroup := groupTotal - prefixGroup[x]
			total := costLonger + costRemoveShorter + costRemoveGroup
			if total < ans {
				ans = total
			}
		}
		for _, v := range g.costs {
			freq[v]++
		}
		prefixCost += groupTotal
		countShorter += c
	}

	if ans < 0 {
		ans = 0
	}
	fmt.Println(ans)
}
