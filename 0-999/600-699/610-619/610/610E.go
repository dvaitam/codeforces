package main

import (
	"bufio"
	"fmt"
	"os"
)

const KMAX = 10

type Node struct {
	cnt   [KMAX][KMAX]int32
	first int8
	last  int8
	len   int32
	lazy  int8
}

var (
	tree []Node
	s    []int8
	k    int
)

func build(idx, l, r int) {
	tree[idx].lazy = -1
	tree[idx].len = int32(r - l + 1)
	if l == r {
		ch := s[l]
		tree[idx].first = ch
		tree[idx].last = ch
		return
	}
	mid := (l + r) / 2
	build(idx*2, l, mid)
	build(idx*2+1, mid+1, r)
	pull(idx)
}

func pull(idx int) {
	left := tree[idx*2]
	right := tree[idx*2+1]
	t := &tree[idx]
	t.first = left.first
	t.last = right.last
	t.len = left.len + right.len
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			t.cnt[i][j] = left.cnt[i][j] + right.cnt[i][j]
		}
	}
	t.cnt[left.last][right.first] += 1
}

func apply(idx int, ch int8) {
	t := &tree[idx]
	t.first = ch
	t.last = ch
	t.lazy = ch
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			t.cnt[i][j] = 0
		}
	}
	if t.len > 1 {
		t.cnt[ch][ch] = t.len - 1
	}
}

func push(idx int) {
	if tree[idx].lazy != -1 {
		apply(idx*2, tree[idx].lazy)
		apply(idx*2+1, tree[idx].lazy)
		tree[idx].lazy = -1
	}
}

func update(idx, l, r, ql, qr int, ch int8) {
	if ql <= l && r <= qr {
		apply(idx, ch)
		return
	}
	push(idx)
	mid := (l + r) / 2
	if ql <= mid {
		update(idx*2, l, mid, ql, qr, ch)
	}
	if qr > mid {
		update(idx*2+1, mid+1, r, ql, qr, ch)
	}
	pull(idx)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(in, &n, &m, &k)
	var str string
	fmt.Fscan(in, &str)
	s = make([]int8, n)
	for i := 0; i < n; i++ {
		s[i] = int8(str[i] - 'a')
	}
	tree = make([]Node, 4*n)
	build(1, 0, n-1)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; m > 0; m-- {
		var tp int
		fmt.Fscan(in, &tp)
		if tp == 1 {
			var l, r int
			var c string
			fmt.Fscan(in, &l, &r, &c)
			update(1, 0, n-1, l-1, r-1, int8(c[0]-'a'))
		} else {
			var perm string
			fmt.Fscan(in, &perm)
			pos := make([]int, k)
			for i := 0; i < k; i++ {
				pos[int(perm[i]-'a')] = i
			}
			root := tree[1]
			ans := 1
			var lastIdx int
			for j := 0; j < k; j++ {
				if pos[j] == k-1 {
					lastIdx = j
					break
				}
			}
			for j := 0; j < k; j++ {
				ans += int(root.cnt[lastIdx][j])
			}
			for i := 0; i < k; i++ {
				if pos[i] == k-1 {
					continue
				}
				for j := 0; j < k; j++ {
					if pos[j] <= pos[i] {
						ans += int(root.cnt[i][j])
					}
				}
			}
			fmt.Fprintln(out, ans)
		}
	}
}
