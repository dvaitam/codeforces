package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type edge struct {
	to, rev int
	cap     int
}

type dinic struct {
	g     [][]edge
	level []int
	iter  []int
}

func newDinic(n int) *dinic {
	g := make([][]edge, n)
	level := make([]int, n)
	iter := make([]int, n)
	return &dinic{g, level, iter}
}

func (d *dinic) addEdge(fr, to, cap int) {
	f := edge{to, len(d.g[to]), cap}
	b := edge{fr, len(d.g[fr]), 0}
	d.g[fr] = append(d.g[fr], f)
	d.g[to] = append(d.g[to], b)
}

func (d *dinic) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	queue := []int{s}
	d.level[s] = 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		for _, e := range d.g[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func (d *dinic) dfs(v, t, f int) int {
	if v == t {
		return f
	}
	for ; d.iter[v] < len(d.g[v]); d.iter[v]++ {
		i := d.iter[v]
		e := d.g[v][i]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			ret := d.dfs(e.to, t, min(f, e.cap))
			if ret > 0 {
				d.g[v][i].cap -= ret
				re := &d.g[e.to][e.rev]
				re.cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *dinic) maxFlow(s, t int) int {
	flow := 0
	for d.bfs(s, t) {
		for i := range d.iter {
			d.iter[i] = 0
		}
		for {
			f := d.dfs(s, t, 1<<30)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]uint64, n)
		bitPresent := make([]bool, 60)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			val := arr[i]
			for val > 0 {
				b := bits.TrailingZeros64(val)
				if b < 60 {
					bitPresent[b] = true
				}
				val &= val - 1
			}
		}

		bitID := make([]int, 60)
		for i := range bitID {
			bitID[i] = -1
		}
		bitCnt := 0
		for b := 0; b < 60; b++ {
			if bitPresent[b] {
				bitID[b] = bitCnt
				bitCnt++
			}
		}

		source := 0
		elemStart := 1
		bitStart := elemStart + n
		sink := bitStart + bitCnt
		flowNet := newDinic(sink + 1)
		INF := n + bitCnt + 5

		for i := 0; i < n; i++ {
			flowNet.addEdge(source, elemStart+i, 1)
			val := arr[i]
			for val > 0 {
				b := bits.TrailingZeros64(val)
				id := bitID[b]
				if id >= 0 {
					flowNet.addEdge(elemStart+i, bitStart+id, INF)
				}
				val &= val - 1
			}
		}

		for b := 0; b < 60; b++ {
			if bitID[b] >= 0 {
				flowNet.addEdge(bitStart+bitID[b], sink, 1)
			}
		}

		maxFlow := flowNet.maxFlow(source, sink)
		ans := n - maxFlow
		fmt.Fprintln(out, ans)
	}
}
