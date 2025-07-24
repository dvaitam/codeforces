package main

import (
	"bufio"
	"fmt"
	"os"
)

const Base uint64 = 911382323
const MaxL = 1000005

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	pow := make([]uint64, MaxL)
	pow[0] = 1
	for i := 1; i < MaxL; i++ {
		pow[i] = pow[i-1] * Base
	}

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		words := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &words[i])
		}
		prefixMap := make(map[uint64]int)
		for _, w := range words {
			var h uint64
			for j := 0; j < len(w); j++ {
				h = h*Base + uint64(w[j]-'a'+1)
				if _, ok := prefixMap[h]; !ok {
					prefixMap[h] = len(prefixMap)
				}
			}
		}
		m := len(prefixMap)
		edges := make([][]int, m)
		added := make([]bool, m)
		dp0 := make([]int, m)
		dp1 := make([]int, m)

		for _, w := range words {
			L := len(w)
			pref := make([]uint64, L+1)
			for i := 0; i < L; i++ {
				pref[i+1] = pref[i]*Base + uint64(w[i]-'a'+1)
			}
			for i := 1; i <= L; i++ {
				h := pref[i]
				idx := prefixMap[h]
				if !added[idx] {
					if i > 1 {
						ph := pref[i] - pref[1]*pow[i-1]
						if pid, ok := prefixMap[ph]; ok {
							edges[idx] = append(edges[idx], pid)
							edges[pid] = append(edges[pid], idx)
						}
					}
					added[idx] = true
				}
			}
		}

		visited := make([]bool, m)
		type node struct{ u, parent, stage int }
		stack := make([]node, 0)
		ans := 0
		for i := 0; i < m; i++ {
			if visited[i] {
				continue
			}
			stack = append(stack, node{u: i, parent: -1, stage: 0})
			for len(stack) > 0 {
				p := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if p.stage == 0 {
					visited[p.u] = true
					stack = append(stack, node{u: p.u, parent: p.parent, stage: 1})
					for _, v := range edges[p.u] {
						if v != p.parent {
							stack = append(stack, node{u: v, parent: p.u, stage: 0})
						}
					}
				} else {
					pick := 1
					notpick := 0
					for _, v := range edges[p.u] {
						if v == p.parent {
							continue
						}
						if dp0[v] > dp1[v] {
							notpick += dp0[v]
						} else {
							notpick += dp1[v]
						}
						pick += dp0[v]
					}
					dp0[p.u] = notpick
					dp1[p.u] = pick
					if p.parent == -1 {
						if dp0[p.u] > dp1[p.u] {
							ans += dp0[p.u]
						} else {
							ans += dp1[p.u]
						}
					}
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
