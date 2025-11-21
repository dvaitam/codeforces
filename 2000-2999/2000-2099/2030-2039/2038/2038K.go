package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type modification struct {
	node int32
	prev int32
}

var (
	xs      []int64
	tree    []int32
	mods    []modification
	hVals   []int64
	nCoords int
)

const infValue int64 = 1 << 62

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func lineValue(idx int32, x int64) int64 {
	if idx < 0 {
		return infValue
	}
	s := int64(idx)
	return s*x + hVals[idx+1]
}

func better(idx1, idx2 int32, x int64) bool {
	if idx2 < 0 {
		return true
	}
	v1 := lineValue(idx1, x)
	v2 := lineValue(idx2, x)
	if v1 == v2 {
		return idx1 < idx2
	}
	return v1 < v2
}

func record(node int, prev int32) {
	mods = append(mods, modification{int32(node), prev})
}

func addLine(node, l, r int, newIdx int32) {
	cur := tree[node]
	if cur < 0 {
		record(node, cur)
		tree[node] = newIdx
		return
	}
	leftX := xs[l]
	mid := (l + r) >> 1
	midX := xs[mid]
	best := cur
	if better(newIdx, best, midX) {
		record(node, tree[node])
		tree[node] = newIdx
		newIdx = best
		best = tree[node]
	}
	if l == r {
		return
	}
	leftBetter := better(newIdx, best, leftX)
	midBetter := better(newIdx, best, midX)
	if leftBetter != midBetter {
		addLine(node<<1, l, mid, newIdx)
	} else {
		addLine(node<<1|1, mid+1, r, newIdx)
	}
}

func query(node, l, r, idx int) int32 {
	res := tree[node]
	if l == r {
		return res
	}
	mid := (l + r) >> 1
	var other int32
	if idx <= mid {
		other = query(node<<1, l, mid, idx)
	} else {
		other = query(node<<1|1, mid+1, r, idx)
	}
	if better(other, res, xs[idx]) {
		return other
	}
	return res
}

func rollback(target int) {
	for len(mods) > target {
		last := mods[len(mods)-1]
		mods = mods[:len(mods)-1]
		tree[last.node] = last.prev
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}

	if n == 1 {
		g1 := int64(gcd(1, a))
		h1 := int64(gcd(1, b))
		fmt.Fprintln(out, g1+h1)
		return
	}

	gVals := make([]int64, n+2)
	hVals = make([]int64, n+2)
	for i := 1; i <= n; i++ {
		gVals[i] = int64(gcd(i, a))
		hVals[i] = int64(gcd(i, b))
	}

	w := make([]int64, n)
	for i := 1; i <= n-1; i++ {
		w[i] = gVals[i] - gVals[i+1]
	}

	coords := make([]int64, n-1)
	for i := 1; i <= n-1; i++ {
		coords[i-1] = w[i]
	}
	sort.Slice(coords, func(i, j int) bool { return coords[i] < coords[j] })
	xs = xs[:0]
	xs = make([]int64, 0, len(coords))
	for _, v := range coords {
		if len(xs) == 0 || xs[len(xs)-1] != v {
			xs = append(xs, v)
		}
	}
	nCoords = len(xs)
	idxMap := make(map[int64]int, nCoords)
	for i, v := range xs {
		idxMap[v] = i
	}

	treeSize := 4 * nCoords
	tree = make([]int32, treeSize)
	for i := range tree {
		tree[i] = -1
	}
	mods = make([]modification, 0)

	snapBefore := make([]int, n)
	for s := 0; s <= n-1; s++ {
		snapBefore[s] = len(mods)
		addLine(1, 0, nCoords-1, int32(s))
	}

	sVals := make([]int64, n+1)
	sVals[n] = int64(n - 1)

	wIdx := make([]int, n)
	for i := 1; i <= n-1; i++ {
		wIdx[i] = idxMap[w[i]]
	}

	currMax := n - 1
	for i := n - 1; i >= 1; i-- {
		target := int(sVals[i+1])
		for currMax > target {
			rollback(snapBefore[currMax])
			currMax--
		}
		idx := wIdx[i]
		best := query(1, 0, nCoords-1, idx)
		if best < 0 {
			best = 0
		}
		sVals[i] = int64(best)
	}

	var gSum int64
	for i := 1; i <= n; i++ {
		gSum += gVals[i]
	}
	var hSum int64
	for i := 1; i <= n; i++ {
		hSum += hVals[i]
	}
	total := gSum + hSum + gVals[n]*int64(n-1)
	for i := 1; i <= n-1; i++ {
		total += w[i]*sVals[i] + hVals[sVals[i]+1]
	}

	fmt.Fprintln(out, total)
}
