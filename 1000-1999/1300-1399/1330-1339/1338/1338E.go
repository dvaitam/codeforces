package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)

	// number of 64-bit words per row
	words := (n + 63) >> 6
	adj := make([][]uint64, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]uint64, words)
		var line string
		fmt.Fscan(reader, &line)
		idx := 0
		for _, c := range line {
			var v int
			if c >= '0' && c <= '9' {
				v = int(c - '0')
			} else {
				v = int(c-'A') + 10
			}
			for b := 0; b < 4 && idx < n; b++ {
				if v&(1<<(3-b)) != 0 {
					adj[i][idx>>6] |= 1 << (uint(idx) & 63)
				}
				idx++
			}
		}
	}

	// Tarjan's algorithm for SCC
	index := 0
	stack := make([]int, 0, n)
	inStack := make([]bool, n)
	dfn := make([]int, n)
	low := make([]int, n)
	compID := make([]int, n)
	compCnt := 0

	var dfs func(int)
	dfs = func(v int) {
		index++
		dfn[v] = index
		low[v] = index
		stack = append(stack, v)
		inStack[v] = true
		for to := 0; to < n; to++ {
			if v == to {
				continue
			}
			if (adj[v][to>>6]>>(uint(to)&63))&1 == 0 {
				continue
			}
			if dfn[to] == 0 {
				dfs(to)
				if low[to] < low[v] {
					low[v] = low[to]
				}
			} else if inStack[to] && dfn[to] < low[v] {
				low[v] = dfn[to]
			}
		}
		if low[v] == dfn[v] {
			for {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				inStack[u] = false
				compID[u] = compCnt
				if u == v {
					break
				}
			}
			compCnt++
		}
	}

	for i := 0; i < n; i++ {
		if dfn[i] == 0 {
			dfs(i)
		}
	}

	comps := make([]struct{ rep, sz int }, compCnt)
	for i := 0; i < n; i++ {
		id := compID[i]
		if comps[id].sz == 0 {
			comps[id].rep = i
		}
		comps[id].sz++
	}

	// sort components by orientation between representatives
	// since condensation is a transitive tournament
	sort.Slice(comps, func(i, j int) bool {
		u := comps[i].rep
		v := comps[j].rep
		return (adj[u][v>>6]>>(uint(v)&63))&1 != 0
	})

	INF := int64(614 * n)
	prefix := 0
	var ans int64
	for _, c := range comps {
		sz := c.sz
		if sz == 3 {
			ans += 9
		}
		if prefix > 0 {
			p := int64(prefix)
			s := int64(sz)
			ans += p*s + s*p*INF
		}
		prefix += sz
	}

	fmt.Fprintln(writer, ans)
}
