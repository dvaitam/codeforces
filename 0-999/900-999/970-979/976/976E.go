package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Creature struct {
	hp   int64
	dmg  int64
	diff int64
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}
	creatures := make([]Creature, n)
	base := int64(0)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &creatures[i].hp, &creatures[i].dmg)
		base += creatures[i].dmg
		d := creatures[i].hp - creatures[i].dmg
		if d < 0 {
			d = 0
		}
		creatures[i].diff = d
	}

	order := make([]int, n)
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		return creatures[order[i]].diff > creatures[order[j]].diff
	})

	pos := make([]int, n)
	arr := make([]int64, n)
	for idx, id := range order {
		pos[id] = idx
		arr[idx] = creatures[id].diff
	}

	pre := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pre[i+1] = pre[i] + arr[i]
	}

	if b > n {
		b = n
	}

	baseline := base + pre[b]
	if b == 0 {
		fmt.Fprintln(out, base)
		return
	}

	ans := baseline
	threshold := arr[b-1]
	for i := 0; i < n; i++ {
		inc := (creatures[i].hp << uint(a)) - creatures[i].dmg
		if inc <= 0 {
			continue
		}
		var cand int64
		if pos[i] < b {
			cand = baseline + inc - creatures[i].diff
		} else {
			cand = baseline
			if inc > threshold {
				cand = cand - threshold + inc
			}
		}
		if cand > ans {
			ans = cand
		}
	}

	if ans < base {
		ans = base
	}
	fmt.Fprintln(out, ans)
}
