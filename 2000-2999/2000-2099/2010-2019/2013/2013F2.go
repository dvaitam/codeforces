package main

import (
	"bufio"
	"fmt"
	"os"
)

type edgeVal struct {
	to  int
	val int
}

type topVal struct {
	val int
	to  int
}

type sparseTable struct {
	log  []int
	st   [][]int
	minQ bool
}

func newSparse(arr []int, useMin bool) *sparseTable {
	n := len(arr)
	log := make([]int, n+1)
	for i := 2; i <= n; i++ {
		log[i] = log[i/2] + 1
	}
	k := log[n] + 1
	st := make([][]int, k)
	st[0] = append([]int(nil), arr...)
	for j := 1; j < k; j++ {
		l := 1 << (j - 1)
		st[j] = make([]int, n-(1<<j)+1)
		for i := 0; i+(1<<j) <= n; i++ {
			if useMin {
				a := st[j-1][i]
				b := st[j-1][i+l]
				if a < b {
					st[j][i] = a
				} else {
					st[j][i] = b
				}
			} else {
				a := st[j-1][i]
				b := st[j-1][i+l]
				if a > b {
					st[j][i] = a
				} else {
					st[j][i] = b
				}
			}
		}
	}
	return &sparseTable{log: log, st: st, minQ: useMin}
}

func (s *sparseTable) query(l, r int) int {
	if l > r {
		if s.minQ {
			return 1 << 30
		}
		return -1 << 30
	}
	k := s.log[r-l+1]
	if s.minQ {
		a := s.st[k][l]
		b := s.st[k][r-(1<<k)+1]
		if a < b {
			return a
		}
		return b
	}
	a := s.st[k][l]
	b := s.st[k][r-(1<<k)+1]
	if a > b {
		return a
	}
	return b
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
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			var a, b int
			fmt.Fscan(in, &a, &b)
			a--
			b--
			adj[a] = append(adj[a], b)
			adj[b] = append(adj[b], a)
		}
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--

		parent := make([]int, n)
		depth := make([]int, n)
		stack := []int{0}
		parent[0] = -1
		for len(stack) > 0 {
			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, nb := range adj[x] {
				if nb == parent[x] {
					continue
				}
				parent[nb] = x
				depth[nb] = depth[x] + 1
				stack = append(stack, nb)
			}
		}

		down := make([]int, n)
		post := make([]int, 0, n)
		stack = append(stack, 0)
		for len(stack) > 0 {
			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			post = append(post, x)
			for _, nb := range adj[x] {
				if nb == parent[x] {
					continue
				}
				stack = append(stack, nb)
			}
		}
		for i := len(post) - 1; i >= 0; i-- {
			x := post[i]
			best := 0
			for _, nb := range adj[x] {
				if nb == parent[x] {
					continue
				}
				if down[nb]+1 > best {
					best = down[nb] + 1
				}
			}
			down[x] = best
		}

		up := make([]int, n)
		var dfsUp func(int)
		dfsUp = func(x int) {
			var best1, best2 int
			var child1 int = -1
			for _, nb := range adj[x] {
				if nb == parent[x] {
					continue
				}
				val := down[nb] + 1
				if val > best1 {
					best2 = best1
					best1 = val
					child1 = nb
				} else if val > best2 {
					best2 = val
				}
			}
			for _, nb := range adj[x] {
				if nb == parent[x] {
					continue
				}
				use := best1
				if nb == child1 {
					use = best2
				}
				if up[x] > use {
					use = up[x]
				}
				up[nb] = use + 1
				dfsUp(nb)
			}
		}
		dfsUp(0)

		top := make([][3]topVal, n)
		for i := 0; i < n; i++ {
			for _, nb := range adj[i] {
				val := 1
				if nb == parent[i] {
					val += up[i]
				} else {
					val += down[nb]
				}
				insertTop(&top[i], topVal{val, nb})
			}
		}

		path := buildPath(u, v, parent)
		wIdx := 0
		bestDepth := 1 << 30
		for idx, node := range path {
			if depth[node] < bestDepth {
				bestDepth = depth[node]
				wIdx = idx
			}
		}

		prev := make([]int, len(path))
		next := make([]int, len(path))
		for i := range path {
			if i == 0 {
				prev[i] = -1
			} else {
				prev[i] = path[i-1]
			}
			if i+1 == len(path) {
				next[i] = -1
			} else {
				next[i] = path[i+1]
			}
		}

		midVal := make([]int, len(path))
		skipL := make([]int, len(path))
		skipR := make([]int, len(path))
		for i, node := range path {
			midVal[i] = bestExcl(top[node], prev[i], next[i])
			skipL[i] = bestExcl(top[node], prev[i], -1)
			skipR[i] = bestExcl(top[node], next[i], -1)
		}

		leftAns := make([]bool, wIdx)
		if wIdx > 0 {
			secSuffix := make([]int, wIdx+1)
			secSuffix[wIdx] = skipL[wIdx]
			for i := wIdx - 1; i >= 0; i-- {
				cur := (wIdx - i) + midVal[i]
				if secSuffix[i+1] > cur {
					secSuffix[i] = secSuffix[i+1]
				} else {
					secSuffix[i] = cur
				}
			}

			D := make([]int, wIdx)
			for i := 0; i < wIdx; i++ {
				D[i] = i + midVal[i] - secSuffix[i+1]
			}
			stMax := newSparse(D, false)
			for i := 0; i < wIdx; i++ {
				mid := (i + wIdx) / 2
				cond1 := skipR[i] > secSuffix[i+1]
				cond2 := false
				if i+1 <= mid {
					if stMax.query(i+1, mid) > i {
						cond2 = true
					}
				}
				leftAns[i] = cond1 || cond2
			}
		}

		rightAns := make([]bool, len(path)-wIdx-1)
		if len(rightAns) > 0 {
			prefix := make([]int, len(path))
			prefix[wIdx] = skipR[wIdx]
			for i := wIdx + 1; i < len(path); i++ {
				cur := (i - wIdx) + midVal[i]
				if cur > prefix[i-1] {
					prefix[i] = cur
				} else {
					prefix[i] = prefix[i-1]
				}
			}

			C := make([]int, len(path))
			for i := wIdx + 1; i < len(path); i++ {
				C[i] = i + prefix[i-1] - midVal[i]
			}
			stMin := newSparse(C[wIdx+1:], true)
			for idx := 0; idx < len(rightAns); idx++ {
				i := wIdx + 1 + idx
				mid := (i + wIdx + 1) / 2
				cond1 := skipL[i] > prefix[i-1]
				cond2 := false
				if mid <= i-1 {
					if i > stMin.query(mid-(wIdx+1), i-1-(wIdx+1)) {
						cond2 = true
					}
				}
				rightAns[idx] = cond1 || cond2
			}
		}

		// Bob starts on each path node, Alice starts at 1 (outside path).
		// For nodes left of wIdx Bob effectively moves first on that segment; Alice wins if leftAns is false.
		// For nodes right of wIdx similar with rightAns.
		for idx, node := range path {
			var aliceWin bool
			if idx < wIdx {
				aliceWin = !leftAns[idx]
			} else if idx > wIdx {
				aliceWin = !rightAns[idx-wIdx-1]
			} else {
				// Bob at wIdx: compute directly on path from 1 to wIdx
				if node == 0 {
					aliceWin = false
				} else {
					p := buildPath(0, node, parent)
					prev2 := make([]int, len(p))
					next2 := make([]int, len(p))
					for i := range p {
						if i == 0 {
							prev2[i] = -1
						} else {
							prev2[i] = p[i-1]
						}
						if i+1 == len(p) {
							next2[i] = -1
						} else {
							next2[i] = p[i+1]
						}
					}
					mid2 := make([]int, len(p))
					sl2 := make([]int, len(p))
					sr2 := make([]int, len(p))
					for i, nd := range p {
						mid2[i] = bestExcl(top[nd], prev2[i], next2[i])
						sl2[i] = bestExcl(top[nd], prev2[i], -1)
						sr2[i] = bestExcl(top[nd], next2[i], -1)
					}
					aliceWin = firstWinsPath(mid2, sl2, sr2)
				}
			}
			if aliceWin {
				fmt.Fprintln(out, "Alice")
			} else {
				fmt.Fprintln(out, "Bob")
			}
		}
	}
}

func insertTop(arr *[3]topVal, tv topVal) {
	for i := 0; i < 3; i++ {
		if tv.val > (*arr)[i].val {
			for j := 2; j > i; j-- {
				(*arr)[j] = (*arr)[j-1]
			}
			(*arr)[i] = tv
			break
		}
	}
}

func bestExcl(tp [3]topVal, a, b int) int {
	for i := 0; i < 3; i++ {
		if tp[i].to == a || tp[i].to == b {
			continue
		}
		return tp[i].val
	}
	return 0
}

func buildPath(u, v int, parent []int) []int {
	seen := make(map[int]bool)
	pathU := make([]int, 0)
	for x := u; x != -1; x = parent[x] {
		seen[x] = true
		pathU = append(pathU, x)
	}
	pathV := make([]int, 0)
	lca := -1
	for x := v; ; x = parent[x] {
		if seen[x] {
			lca = x
			pathV = append(pathV, x)
			break
		}
		pathV = append(pathV, x)
	}
	path := make([]int, 0, len(pathU)+len(pathV))
	for _, x := range pathU {
		path = append(path, x)
		if x == lca {
			break
		}
	}
	for i := len(pathV) - 2; i >= 0; i-- {
		path = append(path, pathV[i])
	}
	return path
}

func firstWinsPath(midVals, skipLeft, skipRight []int) bool {
	k := len(midVals) - 1
	if k == 0 {
		return false
	}
	secSuffix := make([]int, k+1)
	secSuffix[k] = skipLeft[k]
	for i := k - 1; i >= 0; i-- {
		cur := (k - i) + midVals[i]
		if secSuffix[i+1] > cur {
			secSuffix[i] = secSuffix[i+1]
		} else {
			secSuffix[i] = cur
		}
	}

	for i := 0; i <= k/2; i++ {
		var val int
		if i == 0 {
			val = skipRight[0]
		} else {
			val = i + midVals[i]
		}
		if val > secSuffix[i+1] {
			return true
		}
	}
	return false
}
