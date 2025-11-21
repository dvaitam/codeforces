package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	n, k   int
	pos    [][]int
	parent []int
)

func setBit(bs []uint64, idx int) {
	bs[idx>>6] |= 1 << uint(idx&63)
}

func clearBit(bs []uint64, idx int) {
	bs[idx>>6] &^= 1 << uint(idx&63)
}

func copyBitset(src []uint64) []uint64 {
	dst := make([]uint64, len(src))
	copy(dst, src)
	return dst
}

func equalBits(a, b []uint64) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isZero(bs []uint64) bool {
	for _, v := range bs {
		if v != 0 {
			return false
		}
	}
	return true
}

func bitsetToSlice(bs []uint64, nodes []int) []int {
	res := make([]int, 0)
	for i, node := range nodes {
		if (bs[i>>6]>>uint(i&63))&1 == 1 {
			res = append(res, node)
		}
	}
	return res
}

func solve(nodes []int) int {
	m := len(nodes)
	if m == 1 {
		return nodes[0]
	}
	words := (m + 63) >> 6
	idxMap := make(map[int]int, m)
	for i, node := range nodes {
		idxMap[node] = i
	}

	orderIdx := make([][]int, k)
	pref := make([][][]uint64, k)

	for t := 0; t < k; t++ {
		orderIdx[t] = make([]int, m)
		pref[t] = make([][]uint64, m+1)
		data := make([]uint64, (m+1)*words)
		for i := 0; i <= m; i++ {
			pref[t][i] = data[i*words : (i+1)*words]
		}
		order := append([]int(nil), nodes...)
		sort.Slice(order, func(i, j int) bool {
			return pos[t][order[i]] < pos[t][order[j]]
		})
		for i := 0; i < m; i++ {
			copy(pref[t][i+1], pref[t][i])
			localIdx := idxMap[order[i]]
			setBit(pref[t][i+1], localIdx)
			orderIdx[t][localIdx] = i
		}
	}

	fullAll := make([]uint64, words)
	for i := 0; i < m; i++ {
		setBit(fullAll, i)
	}

	for _, r := range nodes {
		ridx := idxMap[r]
		full := copyBitset(fullAll)
		clearBit(full, ridx)
		base := copyBitset(pref[0][orderIdx[0][ridx]])
		if isZero(base) || equalBits(base, full) {
			continue
		}
		comp := make([]uint64, words)
		for i := 0; i < words; i++ {
			comp[i] = full[i] &^ base[i]
		}
		if isZero(comp) {
			continue
		}

		valid := true
		for t := 0; t < k && valid; t++ {
			cur := pref[t][orderIdx[t][ridx]]
			if !equalBits(cur, base) && !equalBits(cur, comp) {
				valid = false
			}
		}
		if !valid {
			continue
		}

		left := bitsetToSlice(base, nodes)
		right := bitsetToSlice(comp, nodes)
		leftRoot := solve(left)
		rightRoot := solve(right)
		parent[leftRoot] = r
		parent[rightRoot] = r
		return r
	}

	panic("root not found")
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n, &k)
	pos = make([][]int, k)
	for i := 0; i < k; i++ {
		pos[i] = make([]int, n+1)
		order := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &order[j])
			pos[i][order[j]] = j
		}
	}

	parent = make([]int, n+1)
	nodes := make([]int, n)
	for i := 0; i < n; i++ {
		nodes[i] = i + 1
	}
	root := solve(nodes)
	parent[root] = -1

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, parent[i])
	}
	fmt.Fprintln(out)
}
