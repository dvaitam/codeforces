package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	const limit = 200000
	first := make([]int, limit+1)
	second := make([]int, limit+1)
	count := make([]int, limit+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		v := arr[i]
		count[v]++
		if count[v] == 1 {
			first[v] = i + 1
		} else if count[v] == 2 {
			second[v] = i + 1
		}
	}

	k0 := 0
	for v := 1; v <= limit; v++ {
		if count[v] == 2 {
			k0 = v
		} else {
			break
		}
	}

	type pair struct{ f, s, val int }
	pairs := make([]pair, 0, k0)
	for v := 1; v <= k0; v++ {
		pairs = append(pairs, pair{first[v], second[v], v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].f < pairs[j].f })

	curSecond := 0
	seen := 0
	maxValSeen := 0
	k := 0
	ok := true
	for _, p := range pairs {
		if p.s <= curSecond {
			ok = false
		}
		curSecond = p.s
		seen++
		if p.val > maxValSeen {
			maxValSeen = p.val
		}
		if ok && seen == maxValSeen {
			k = maxValSeen
		}
	}

	ans := make([]byte, n)
	for i := range ans {
		ans[i] = 'B'
	}
	for v := 1; v <= k; v++ {
		ans[first[v]-1] = 'R'
		ans[second[v]-1] = 'G'
	}
	fmt.Println(string(ans))
}
