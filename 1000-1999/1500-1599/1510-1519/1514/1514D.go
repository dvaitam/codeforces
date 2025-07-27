package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	val int
	cnt int
}

var a []int
var tree []Node

func merge(left, right Node) Node {
	if left.val == right.val {
		return Node{left.val, left.cnt + right.cnt}
	}
	if left.cnt > right.cnt {
		return Node{left.val, left.cnt - right.cnt}
	}
	if right.cnt > left.cnt {
		return Node{right.val, right.cnt - left.cnt}
	}
	return Node{0, 0}
}

func build(idx, l, r int) {
	if l == r {
		tree[idx] = Node{a[l], 1}
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	tree[idx] = merge(tree[idx*2], tree[idx*2+1])
}

func query(idx, l, r, L, R int) Node {
	if L <= l && r <= R {
		return tree[idx]
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

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	a = make([]int, n+1)
	pos := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
		pos[a[i]] = append(pos[a[i]], i)
	}

	tree = make([]Node, 4*(n+1))
	build(1, 1, n)

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		cand := query(1, 1, n, l, r).val
		freq := 0
		if cand != 0 {
			arr := pos[cand]
			left := sort.SearchInts(arr, l)
			right := sort.SearchInts(arr, r+1)
			freq = right - left
		}
		length := r - l + 1
		ans := 1
		if tmp := 2*freq - length; tmp > 1 {
			ans = tmp
		}
		fmt.Fprintln(writer, ans)
	}
}
