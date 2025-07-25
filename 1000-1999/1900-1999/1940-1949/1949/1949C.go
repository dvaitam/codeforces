package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func insertSorted(a []int, x int) []int {
	i := sort.SearchInts(a, x)
	a = append(a, 0)
	copy(a[i+1:], a[i:])
	a[i] = x
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	parent := make([]int, n)
	order := make([]int, 0, n)
	st := []int{0}
	parent[0] = -1
	for len(st) > 0 {
		v := st[len(st)-1]
		st = st[:len(st)-1]
		order = append(order, v)
		for _, to := range g[v] {
			if to == parent[v] {
				continue
			}
			parent[to] = v
			st = append(st, to)
		}
	}

	childVals := make([][]int, n)
	sub := make([]int, n)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		vals := make([]int, 0, len(g[u]))
		for _, v := range g[u] {
			if v == parent[u] {
				continue
			}
			vals = append(vals, sub[v])
		}
		sort.Ints(vals)
		cur := 1
		ok := true
		for _, x := range vals {
			if x == 0 || cur < x {
				ok = false
				break
			}
			cur += x
		}
		if ok {
			sub[u] = cur
		} else {
			sub[u] = 0
		}
		childVals[u] = vals
	}

	type node struct{ u, p, val int }
	st2 := []node{{0, -1, 0}}
	ans := false

	for len(st2) > 0 {
		curNode := st2[len(st2)-1]
		st2 = st2[:len(st2)-1]
		u, p, fromParent := curNode.u, curNode.p, curNode.val

		arr := append([]int(nil), childVals[u]...)
		if p != -1 {
			arr = insertSorted(arr, fromParent)
		}
		sort.Ints(arr)
		m := len(arr)
		preSum := make([]int, m+1)
		for i, x := range arr {
			preSum[i+1] = preSum[i] + x
		}
		diff := make([]int, m)
		for i := range arr {
			diff[i] = arr[i] - preSum[i]
		}
		pre := make([]int, m)
		mx := -1 << 60
		for i := 0; i < m; i++ {
			if diff[i] > mx {
				mx = diff[i]
			}
			pre[i] = mx
		}
		suf := make([]int, m)
		mx = -1 << 60
		for i := m - 1; i >= 0; i-- {
			if diff[i] > mx {
				mx = diff[i]
			}
			suf[i] = mx
		}
		needAll := 0
		if m > 0 {
			needAll = pre[m-1]
			if needAll < 0 {
				needAll = 0
			}
		}
		cur := 1
		ok := true
		if 1 < needAll {
			ok = false
		} else {
			for _, x := range arr {
				if x == 0 || cur < x {
					ok = false
					break
				}
				cur += x
			}
		}
		if ok && cur == n {
			ans = true
		}

		// map for duplicates
		idxMap := make(map[int][]int)
		for i, x := range arr {
			idxMap[x] = append(idxMap[x], i)
		}
		used := make(map[int]int)

		for _, v := range g[u] {
			if v == p {
				continue
			}
			valChild := sub[v]
			idxList := idxMap[valChild]
			idx := idxList[used[valChild]]
			used[valChild]++
			left := -1 << 60
			if idx > 0 {
				left = pre[idx-1]
			}
			right := -1 << 60
			if idx+1 < m {
				right = suf[idx+1] + valChild
			}
			need := left
			if right > need {
				need = right
			}
			if need < 0 {
				need = 0
			}
			var val int
			if 1 < need {
				val = 0
			} else {
				val = 1 + preSum[m] - valChild
			}
			st2 = append(st2, node{v, u, val})
		}
	}

	if ans {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
