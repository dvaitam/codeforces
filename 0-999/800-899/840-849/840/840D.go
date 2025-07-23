package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	cnt int
}

type node struct {
	p [5]pair
	m int
}

func (n *node) add(val, c int) {
	for i := 0; i < n.m; i++ {
		if n.p[i].val == val {
			n.p[i].cnt += c
			return
		}
	}
	if n.m < 5 {
		n.p[n.m] = pair{val: val, cnt: c}
		n.m++
	}
	if n.m > 4 {
		minc := n.p[0].cnt
		for i := 1; i < n.m; i++ {
			if n.p[i].cnt < minc {
				minc = n.p[i].cnt
			}
		}
		j := 0
		for i := 0; i < n.m; i++ {
			n.p[i].cnt -= minc
			if n.p[i].cnt > 0 {
				n.p[j] = n.p[i]
				j++
			}
		}
		n.m = j
	}
}

func merge(a, b node) node {
	var res node
	for i := 0; i < a.m; i++ {
		res.add(a.p[i].val, a.p[i].cnt)
	}
	for i := 0; i < b.m; i++ {
		res.add(b.p[i].val, b.p[i].cnt)
	}
	return res
}

var (
	arr  []int
	seg  []node
	pos  map[int][]int
	n, q int
)

func build(idx, l, r int) {
	if l == r {
		seg[idx].p[0] = pair{val: arr[l], cnt: 1}
		seg[idx].m = 1
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	seg[idx] = merge(seg[idx*2], seg[idx*2+1])
}

func query(idx, l, r, L, R int) node {
	if L <= l && r <= R {
		return seg[idx]
	}
	mid := (l + r) / 2
	if R <= mid {
		return query(idx*2, l, mid, L, R)
	}
	if L > mid {
		return query(idx*2+1, mid+1, r, L, R)
	}
	left := query(idx*2, l, mid, L, R)
	right := query(idx*2+1, mid+1, r, L, R)
	return merge(left, right)
}

func count(val, l, r int) int {
	arr := pos[val]
	if len(arr) == 0 {
		return 0
	}
	left := sort.Search(len(arr), func(i int) bool { return arr[i] >= l })
	right := sort.Search(len(arr), func(i int) bool { return arr[i] > r })
	return right - left
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &q)
	arr = make([]int, n+1)
	pos = make(map[int][]int)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
		pos[arr[i]] = append(pos[arr[i]], i)
	}

	seg = make([]node, 4*n)
	build(1, 1, n)

	for ; q > 0; q-- {
		var l, r, k int
		fmt.Fscan(reader, &l, &r, &k)
		cand := query(1, 1, n, l, r)
		threshold := (r - l + 1) / k
		ans := -1
		for i := 0; i < cand.m; i++ {
			v := cand.p[i].val
			freq := count(v, l, r)
			if freq > threshold {
				if ans == -1 || v < ans {
					ans = v
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
