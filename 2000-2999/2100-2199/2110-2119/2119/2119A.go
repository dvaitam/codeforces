package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD = 1024
const INF int64 = 1 << 60

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var a, b int
		var x, y int64
		fmt.Fscan(in, &a, &b, &x, &y)
		dist := make([]int64, MOD)
		inQueue := make([]bool, MOD)
		for i := 0; i < MOD; i++ {
			dist[i] = INF
		}
		init := a % MOD
		dist[init] = 0
		queue := []int{init}
		inQueue[init] = true
		head := 0
		for head < len(queue) {
			u := queue[head]
			head++
			inQueue[u] = false
			nexts := []int{(u + 1) % MOD, u ^ 1}
			costs := []int64{x, y}
			for idx, v := range nexts {
				cand := dist[u] + costs[idx]
				if cand < dist[v] {
					dist[v] = cand
					if !inQueue[v] {
						queue = append(queue, v)
						inQueue[v] = true
					}
				}
			}
		}
		result := dist[b%MOD]
		if result >= INF/2 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, result)
		}
	}
}
