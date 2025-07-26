package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Player struct {
	a   int
	b   int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		players := make([]Player, n)
		for i := 0; i < n; i++ {
			players[i] = Player{a[i], b[i], i}
		}

		orderA := make([]Player, n)
		orderB := make([]Player, n)
		copy(orderA, players)
		copy(orderB, players)
		sort.Slice(orderA, func(i, j int) bool { return orderA[i].a < orderA[j].a })
		sort.Slice(orderB, func(i, j int) bool { return orderB[i].b < orderB[j].b })

		visited := make([]bool, n)
		queue := make([]int, 0)

		pa := n - 1
		pb := n - 1
		add := func(id int) {
			if !visited[id] {
				visited[id] = true
				queue = append(queue, id)
			}
		}

		add(orderA[pa].idx)
		add(orderB[pb].idx)

		for idx := 0; idx < len(queue); idx++ {
			v := queue[idx]
			for pa >= 0 && orderA[pa].a >= a[v] {
				add(orderA[pa].idx)
				pa--
			}
			for pb >= 0 && orderB[pb].b >= b[v] {
				add(orderB[pb].idx)
				pb--
			}
		}

		res := make([]byte, n)
		for i := 0; i < n; i++ {
			if visited[i] {
				res[i] = '1'
			} else {
				res[i] = '0'
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
