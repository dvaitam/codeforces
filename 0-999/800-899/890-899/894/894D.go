package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	n, m   int
	w      []int // edge weight from parent to node (index by node)
	parent []int
	dist   [][]int   // sorted distances from node to nodes in its subtree
	pref   [][]int64 // prefix sums of distances
)

func build() {
	dist = make([][]int, n+1)
	pref = make([][]int64, n+1)
	// process nodes from n down to 1
	for i := n; i >= 1; i-- {
		left := i * 2
		right := left + 1
		lExist := left <= n
		rExist := right <= n
		if !lExist && !rExist {
			dist[i] = []int{0}
			pref[i] = []int64{0}
			continue
		}
		if lExist && !rExist {
			lvals := dist[left]
			off := w[left]
			arr := make([]int, 1+len(lvals))
			arr[0] = 0
			for j, v := range lvals {
				arr[j+1] = v + off
			}
			pre := make([]int64, len(arr))
			var s int64
			for j, v := range arr {
				s += int64(v)
				pre[j] = s
			}
			dist[i] = arr
			pref[i] = pre
			continue
		}
		if !lExist && rExist {
			rvals := dist[right]
			off := w[right]
			arr := make([]int, 1+len(rvals))
			arr[0] = 0
			for j, v := range rvals {
				arr[j+1] = v + off
			}
			pre := make([]int64, len(arr))
			var s int64
			for j, v := range arr {
				s += int64(v)
				pre[j] = s
			}
			dist[i] = arr
			pref[i] = pre
			continue
		}
		// both children exist
		lvals := dist[left]
		lOff := w[left]
		rvals := dist[right]
		rOff := w[right]
		arr := make([]int, 1+len(lvals)+len(rvals))
		arr[0] = 0
		i1, i2, pos := 0, 0, 1
		for i1 < len(lvals) && i2 < len(rvals) {
			v1 := lvals[i1] + lOff
			v2 := rvals[i2] + rOff
			if v1 < v2 {
				arr[pos] = v1
				i1++
			} else {
				arr[pos] = v2
				i2++
			}
			pos++
		}
		for i1 < len(lvals) {
			arr[pos] = lvals[i1] + lOff
			i1++
			pos++
		}
		for i2 < len(rvals) {
			arr[pos] = rvals[i2] + rOff
			i2++
			pos++
		}
		pre := make([]int64, len(arr))
		var s int64
		for j, v := range arr {
			s += int64(v)
			pre[j] = s
		}
		dist[i] = arr
		pref[i] = pre
	}
}

func calc(node int, limit int) int64 {
	arr := dist[node]
	if len(arr) == 0 {
		return 0
	}
	idx := sort.Search(len(arr), func(i int) bool { return arr[i] > limit })
	if idx == 0 {
		return 0
	}
	return int64(limit)*int64(idx) - pref[node][idx-1]
}

func query(a int, h int) int64 {
	ans := calc(a, h)
	child := a
	anc := parent[child]
	distUp := 0
	for anc != 0 {
		distUp += w[child]
		rem := h - distUp
		if rem <= 0 {
			break
		}
		ans += int64(rem)
		left := anc * 2
		right := left + 1
		if left <= n && left != child {
			r := rem - w[left]
			if r > 0 {
				ans += calc(left, r)
			}
		}
		if right <= n && right != child {
			r := rem - w[right]
			if r > 0 {
				ans += calc(right, r)
			}
		}
		child = anc
		anc = parent[child]
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	w = make([]int, n+1)
	parent = make([]int, n+1)
	parent[1] = 0
	for i := 2; i <= n; i++ {
		var L int
		fmt.Fscan(in, &L)
		w[i] = L
		parent[i] = i / 2
	}

	build()

	for ; m > 0; m-- {
		var a, h int
		fmt.Fscan(in, &a, &h)
		res := query(a, h)
		fmt.Fprintln(out, res)
	}
}
