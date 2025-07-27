package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Player struct {
	p    int
	team int // +1 for first team, -1 for second team
}

func computeDiff(k int, a []int, b []int) int64 {
	n := len(a)
	m := len(b)
	players := make([]Player, 0, n+m)
	for _, v := range a {
		players = append(players, Player{p: v, team: 1})
	}
	for _, v := range b {
		players = append(players, Player{p: v, team: -1})
	}
	sort.Slice(players, func(i, j int) bool {
		if players[i].p == players[j].p {
			return players[i].team > players[j].team
		}
		return players[i].p > players[j].p
	})
	if len(players) == 0 {
		return 0
	}
	device := k
	prev := players[0].p
	var diff int64
	diff += int64(players[0].team) * int64(device)
	for i := 1; i < len(players); i++ {
		d := players[i].p - prev
		device += d
		if device < 0 {
			device = 0
		}
		prev = players[i].p
		diff += int64(players[i].team) * int64(device)
	}
	return diff
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	a := make([]int, n)
	b := make([]int, m)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
	}
	for ; q > 0; q-- {
		var tp int
		fmt.Fscan(in, &tp)
		if tp == 1 {
			var pos, x int
			fmt.Fscan(in, &pos, &x)
			if pos >= 1 && pos <= n {
				a[pos-1] = x
			}
		} else if tp == 2 {
			var pos, x int
			fmt.Fscan(in, &pos, &x)
			if pos >= 1 && pos <= m {
				b[pos-1] = x
			}
		} else if tp == 3 {
			var x int
			fmt.Fscan(in, &x)
			ans := computeDiff(x, a, b)
			fmt.Fprintln(out, ans)
		}
	}
}
