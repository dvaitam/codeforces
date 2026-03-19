package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var A, B string
		fmt.Fscan(reader, &A)
		fmt.Fscan(reader, &B)

		var dir [20][20]bool
		var used [20]bool
		var undir [20][20]bool

		for i := 0; i < n; i++ {
			u := int(A[i] - 'a')
			v := int(B[i] - 'a')
			if u != v {
				dir[u][v] = true
				used[u] = true
				used[v] = true
				undir[u][v] = true
				undir[v][u] = true
			}
		}

		ans := 0
		var vis [20]bool

		for s := 0; s < 20; s++ {
			if !used[s] || vis[s] {
				continue
			}

			// BFS to find connected component
			queue := []int{s}
			vis[s] = true
			comp := []int{}

			for head := 0; head < len(queue); head++ {
				u := queue[head]
				comp = append(comp, u)
				for v := 0; v < 20; v++ {
					if undir[u][v] && !vis[v] {
						vis[v] = true
						queue = append(queue, v)
					}
				}
			}

			m := len(comp)
			if m == 1 {
				continue
			}

			// Map component vertices to local indices
			pos := make([]int, 20)
			for i := range pos {
				pos[i] = -1
			}
			for i, v := range comp {
				pos[v] = i
			}

			// Build local directed edge masks
			outMask := make([]uint32, m)
			inMask := make([]uint32, m)
			for i, u := range comp {
				for v := 0; v < 20; v++ {
					if dir[u][v] {
						j := pos[v]
						if j >= 0 {
							outMask[i] |= 1 << j
							inMask[j] |= 1 << i
						}
					}
				}
			}

			// Find maximum acyclic subgraph (max subset of vertices inducing a DAG)
			limit := 1 << m
			acyclic := make([]uint8, limit)
			acyclic[0] = 1
			best := 0

			for mask := 1; mask < limit; mask++ {
				smask := uint32(mask)
				x := smask
				ok := false
				for x != 0 {
					lsb := x & (-x)
					i := bits.TrailingZeros32(lsb)
					prev := mask ^ int(lsb)
					if acyclic[prev] == 1 {
						// vertex i can be removed from the subset and result is acyclic
						// check that i is a source or sink in the subgraph induced by mask
						if (outMask[i]&smask) == 0 || (inMask[i]&smask) == 0 {
							ok = true
							break
						}
					}
					x -= lsb
				}
				if ok {
					acyclic[mask] = 1
					pc := bits.OnesCount32(smask)
					if pc > best {
						best = pc
					}
				}
			}

			fvs := m - best
			ans += (m - 1) + fvs
		}

		fmt.Fprintln(writer, ans)
	}
}
