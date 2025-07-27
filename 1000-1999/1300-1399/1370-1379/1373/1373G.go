package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxV = 400005
const negInf = -1 << 60

type node struct {
	cnt  int
	best int
}

var seg []node

func combine(l, r node) node {
	cnt := l.cnt + r.cnt
	best := l.best
	if r.best+l.cnt > best {
		best = r.best + l.cnt
	}
	return node{cnt: cnt, best: best}
}

func update(idx, delta, v, l, r int) {
	if l == r {
		seg[idx].cnt += delta
		if seg[idx].cnt > 0 {
			seg[idx].best = l + seg[idx].cnt - 1
		} else {
			seg[idx].best = negInf
		}
		return
	}
	mid := (l + r) >> 1
	if v <= mid {
		update(idx<<1, delta, v, l, mid)
	} else {
		update(idx<<1|1, delta, v, mid+1, r)
	}
	seg[idx] = combine(seg[idx<<1], seg[idx<<1|1])
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k, m int
	fmt.Fscan(reader, &n, &k, &m)

	seg = make([]node, maxV*4)
	board := make(map[[2]int]bool)

	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		p := [2]int{x, y}
		v := y + abs(x-k)
		if board[p] {
			delete(board, p)
			update(1, -1, v, 0, maxV-1)
		} else {
			board[p] = true
			update(1, 1, v, 0, maxV-1)
		}
		r := seg[1].best
		if seg[1].cnt > r {
			r = seg[1].cnt
		}
		ans := r - n
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(writer, ans)
	}
}
