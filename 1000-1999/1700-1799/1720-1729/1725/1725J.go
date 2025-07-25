package main

import (
	"bufio"
	"fmt"
	"os"
)

type Edge struct {
	to int
	w  int64
}

var (
	g          [][]Edge
	parent     []int
	parW       []int64
	subSum     []int64
	downDist   []int64
	secondDist []int64
	diam       []int64
	upDist     []int64
	upDia      []int64
)

func dfs1(u, p int) {
	parent[u] = p
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		dfs1(e.to, u)
		parW[e.to] = e.w
		subSum[u] += subSum[e.to] + e.w
		d := downDist[e.to] + e.w
		if d > downDist[u] {
			secondDist[u] = downDist[u]
			downDist[u] = d
		} else if d > secondDist[u] {
			secondDist[u] = d
		}
		if diam[e.to] > diam[u] {
			diam[u] = diam[e.to]
		}
	}
	if downDist[u]+secondDist[u] > diam[u] {
		diam[u] = downDist[u] + secondDist[u]
	}
}

func dfs2(u, p int) {
	// collect children information
	type childInfo struct {
		v int
		w int64
	}
	children := make([]childInfo, 0)
	for _, e := range g[u] {
		if e.to == p {
			continue
		}
		children = append(children, childInfo{e.to, e.w})
	}
	k := len(children)
	pref1Val := make([]int64, k)
	pref1ID := make([]int, k)
	pref2Val := make([]int64, k)
	pref2ID := make([]int, k)
	prefDia := make([]int64, k)
	var b1Val, b2Val int64 = -1, -1
	var b1ID, b2ID int = -1, -1
	for i := 0; i < k; i++ {
		v := children[i].v
		w := children[i].w
		val := downDist[v] + w
		if val > b1Val {
			b2Val, b2ID = b1Val, b1ID
			b1Val, b1ID = val, i
		} else if val > b2Val {
			b2Val, b2ID = val, i
		}
		if i == 0 {
			pref1Val[i] = val
			pref1ID[i] = i
			pref2Val[i] = -1
			pref2ID[i] = -1
			prefDia[i] = diam[v]
		} else {
			pref1Val[i] = b1Val
			pref1ID[i] = b1ID
			pref2Val[i] = b2Val
			pref2ID[i] = b2ID
			pd := prefDia[i-1]
			if diam[v] > pd {
				pd = diam[v]
			}
			prefDia[i] = pd
		}
	}
	suf1Val := make([]int64, k)
	suf1ID := make([]int, k)
	suf2Val := make([]int64, k)
	suf2ID := make([]int, k)
	sufDia := make([]int64, k)
	b1Val, b2Val = -1, -1
	b1ID, b2ID = -1, -1
	for i := k - 1; i >= 0; i-- {
		v := children[i].v
		w := children[i].w
		val := downDist[v] + w
		if val > b1Val {
			b2Val, b2ID = b1Val, b1ID
			b1Val, b1ID = val, i
		} else if val > b2Val {
			b2Val, b2ID = val, i
		}
		if i == k-1 {
			suf1Val[i] = val
			suf1ID[i] = i
			suf2Val[i] = -1
			suf2ID[i] = -1
			sufDia[i] = diam[v]
		} else {
			suf1Val[i] = b1Val
			suf1ID[i] = b1ID
			suf2Val[i] = b2Val
			suf2ID[i] = b2ID
			sd := sufDia[i+1]
			if diam[v] > sd {
				sd = diam[v]
			}
			sufDia[i] = sd
		}
	}
	for i := 0; i < k; i++ {
		v := children[i].v
		w := children[i].w
		md := upDist[u]
		if i > 0 && pref1Val[i-1] > md {
			md = pref1Val[i-1]
		}
		if i+1 < k && suf1Val[i+1] > md {
			md = suf1Val[i+1]
		}
		upDist[v] = md + w

		option := upDia[u]
		diaEx := int64(0)
		if i > 0 && prefDia[i-1] > diaEx {
			diaEx = prefDia[i-1]
		}
		if i+1 < k && sufDia[i+1] > diaEx {
			diaEx = sufDia[i+1]
		}
		if diaEx > option {
			option = diaEx
		}
		type pair struct {
			val int64
			id  int
		}
		cand := []pair{{upDist[u], -1}}
		if i > 0 {
			cand = append(cand, pair{pref1Val[i-1], pref1ID[i-1]})
			if pref2Val[i-1] >= 0 {
				cand = append(cand, pair{pref2Val[i-1], pref2ID[i-1]})
			}
		}
		if i+1 < k {
			cand = append(cand, pair{suf1Val[i+1], suf1ID[i+1]})
			if suf2Val[i+1] >= 0 {
				cand = append(cand, pair{suf2Val[i+1], suf2ID[i+1]})
			}
		}
		best1 := pair{-1, -1}
		best2 := pair{-1, -1}
		for _, c := range cand {
			if c.id == i || c.val < 0 {
				continue
			}
			if c.val > best1.val {
				best2 = best1
				best1 = c
			} else if c.val > best2.val {
				best2 = c
			}
		}
		cross := best1.val
		if best2.val >= 0 {
			cross += best2.val
		}
		if cross > option {
			option = cross
		}
		upDia[v] = option
		dfs2(v, u)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g = make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		var w int64
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		g[u] = append(g[u], Edge{v, w})
		g[v] = append(g[v], Edge{u, w})
	}

	parent = make([]int, n)
	parW = make([]int64, n)
	subSum = make([]int64, n)
	downDist = make([]int64, n)
	secondDist = make([]int64, n)
	diam = make([]int64, n)
	upDist = make([]int64, n)
	upDia = make([]int64, n)

	dfs1(0, -1)
	dfs2(0, -1)

	totalSum := subSum[0]
	ans := 2*totalSum - diam[0]
	for v := 1; v < n; v++ {
		sumB := subSum[v]
		sumA := totalSum - sumB - parW[v]
		cost := (2*sumB - diam[v]) + (2*sumA - upDia[v])
		if cost < ans {
			ans = cost
		}
	}
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, ans)
	out.Flush()
}
