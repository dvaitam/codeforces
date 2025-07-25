package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf int64 = -1 << 60

var (
	g     [][]int
	size  []int
	n, k  int
	base  int64
	dpSel [][]int64
	dpNot [][]int64
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func dfs1(v, p int) {
	size[v] = 1
	for _, to := range g[v] {
		if to == p {
			continue
		}
		dfs1(to, v)
		s := size[to]
		base += int64(2) * int64(s) * int64(n-s)
		size[v] += s
	}
}

func dfs2(v, p int) {
	dpSel[v] = make([]int64, min(size[v], k)+1)
	dpNot[v] = make([]int64, min(size[v], k)+1)
	for i := 0; i < len(dpSel[v]); i++ {
		dpSel[v][i] = negInf
		dpNot[v][i] = negInf
	}
	dpSel[v][1] = 0
	dpNot[v][0] = 0
	curSize := 1
	for _, to := range g[v] {
		if to == p {
			continue
		}
		dfs2(to, v)
		w := int64(size[to]) * int64(n-size[to])
		newSel := make([]int64, min(curSize+size[to], k)+1)
		newNot := make([]int64, min(curSize+size[to], k)+1)
		for i := 0; i < len(newSel); i++ {
			newSel[i] = negInf
			newNot[i] = negInf
		}
		maxI := min(curSize, k)
		for i := 0; i <= maxI; i++ {
			maxJ := min(size[to], k-i)
			for j := 0; j <= maxJ; j++ {
				if i < len(dpNot[v]) && j < len(dpNot[to]) {
					val := dpNot[v][i]
					if val > negInf {
						valChild := dpNot[to][j]
						if dpSel[to][j] > valChild {
							valChild = dpSel[to][j]
						}
						if valChild > negInf && newNot[i+j] < val+valChild {
							newNot[i+j] = val + valChild
						}
					}
				}
				if i < len(dpSel[v]) {
					if j < len(dpNot[to]) {
						val := dpSel[v][i]
						if val > negInf && dpNot[to][j] > negInf {
							if newSel[i+j] < val+dpNot[to][j] {
								newSel[i+j] = val + dpNot[to][j]
							}
						}
					}
					if j < len(dpSel[to]) {
						val := dpSel[v][i]
						if val > negInf && dpSel[to][j] > negInf {
							if newSel[i+j] < val+dpSel[to][j]+w {
								newSel[i+j] = val + dpSel[to][j] + w
							}
						}
					}
				}
			}
		}
		dpSel[v] = newSel
		dpNot[v] = newNot
		curSize += size[to]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	out := bufio.NewWriter(os.Stdout)
	for ; T > 0; T-- {
		fmt.Fscan(in, &n, &k)
		g = make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		size = make([]int, n)
		base = 0
		dfs1(0, -1)
		dpSel = make([][]int64, n)
		dpNot = make([][]int64, n)
		dfs2(0, -1)
		best := int64(0)
		if k < len(dpSel[0]) && dpSel[0][k] > best {
			best = dpSel[0][k]
		}
		if k < len(dpNot[0]) && dpNot[0][k] > best {
			best = dpNot[0][k]
		}
		res := base - best
		fmt.Fprintln(out, res)
	}
	out.Flush()
}
