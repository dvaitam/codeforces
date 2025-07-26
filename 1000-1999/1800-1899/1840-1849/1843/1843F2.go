package main

import (
	"bufio"
	"fmt"
	"os"
)

const LOG = 20

type Seg struct {
	sum     int
	maxPref int
	minPref int
	maxSuf  int
	minSuf  int
	maxSub  int
	minSub  int
}

func makeSeg(w int) Seg {
	if w > 0 {
		return Seg{sum: w, maxPref: w, minPref: 0, maxSuf: w, minSuf: 0, maxSub: w, minSub: 0}
	}
	return Seg{sum: w, maxPref: 0, minPref: w, maxSuf: 0, minSuf: w, maxSub: 0, minSub: w}
}

func merge(a, b Seg) Seg {
	res := Seg{}
	res.sum = a.sum + b.sum
	if a.maxPref > a.sum+b.maxPref {
		res.maxPref = a.maxPref
	} else {
		res.maxPref = a.sum + b.maxPref
	}
	if a.minPref < a.sum+b.minPref {
		res.minPref = a.minPref
	} else {
		res.minPref = a.sum + b.minPref
	}
	if b.maxSuf > b.sum+a.maxSuf {
		res.maxSuf = b.maxSuf
	} else {
		res.maxSuf = b.sum + a.maxSuf
	}
	if b.minSuf < b.sum+a.minSuf {
		res.minSuf = b.minSuf
	} else {
		res.minSuf = b.sum + a.minSuf
	}
	cross := a.maxSuf + b.maxPref
	res.maxSub = a.maxSub
	if b.maxSub > res.maxSub {
		res.maxSub = b.maxSub
	}
	if cross > res.maxSub {
		res.maxSub = cross
	}
	cross2 := a.minSuf + b.minPref
	res.minSub = a.minSub
	if b.minSub < res.minSub {
		res.minSub = b.minSub
	}
	if cross2 < res.minSub {
		res.minSub = cross2
	}
	return res
}

func reverseSeg(s Seg) Seg {
	return Seg{sum: s.sum, maxPref: s.maxSuf, minPref: s.minSuf, maxSuf: s.maxPref, minSuf: s.minPref, maxSub: s.maxSub, minSub: s.minSub}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		size := n + 5
		parent := make([]int, size)
		weight := make([]int, size)
		depth := make([]int, size)
		up := make([][]int, LOG)
		seg := make([][]Seg, LOG)
		for i := 0; i < LOG; i++ {
			up[i] = make([]int, size)
			seg[i] = make([]Seg, size)
		}
		weight[1] = 1
		cur := 1
		for i := 0; i < n; i++ {
			var op string
			fmt.Fscan(reader, &op)
			if op == "+" {
				var v, x int
				fmt.Fscan(reader, &v, &x)
				cur++
				parent[cur] = v
				weight[cur] = x
				depth[cur] = depth[v] + 1
				up[0][cur] = v
				seg[0][cur] = makeSeg(x)
				for j := 1; j < LOG; j++ {
					up[j][cur] = up[j-1][up[j-1][cur]]
					seg[j][cur] = merge(seg[j-1][cur], seg[j-1][up[j-1][cur]])
				}
			} else {
				var u, v, k int
				fmt.Fscan(reader, &u, &v, &k)
				l := lca(u, v, up, depth)
				s1 := getPathUp(u, l, up, seg, depth)
				s2 := getPathUp(v, l, up, seg, depth)
				s2 = reverseSeg(s2)
				total := merge(merge(s1, makeSeg(weight[l])), s2)
				if k >= total.minSub && k <= total.maxSub {
					fmt.Fprintln(writer, "Yes")
				} else {
					fmt.Fprintln(writer, "No")
				}
			}
		}
	}
}

func lca(u, v int, up [][]int, depth []int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	for i := LOG - 1; i >= 0; i-- {
		if diff&(1<<i) != 0 {
			u = up[i][u]
		}
	}
	if u == v {
		return u
	}
	for i := LOG - 1; i >= 0; i-- {
		if up[i][u] != up[i][v] {
			u = up[i][u]
			v = up[i][v]
		}
	}
	return up[0][u]
}

func getPathUp(u, anc int, up [][]int, seg [][]Seg, depth []int) Seg {
	res := Seg{}
	diff := depth[u] - depth[anc]
	for i := LOG - 1; i >= 0; i-- {
		if diff&(1<<i) != 0 {
			res = merge(res, seg[i][u])
			u = up[i][u]
		}
	}
	return res
}
