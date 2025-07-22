package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Node struct {
	covered bool
	c       [4]*Node
}

func (node *Node) update(x1, x2, y1, y2, xl, xr, yl, yr int) {
	if node.covered || x2 <= xl || xr <= x1 || y2 <= yl || yr <= y1 {
		return
	}
	if x1 <= xl && xr <= x2 && y1 <= yl && yr <= y2 {
		node.covered = true
		node.c = [4]*Node{}
		return
	}
	if xl+1 == xr && yl+1 == yr {
		node.covered = true
		node.c = [4]*Node{}
		return
	}
	mx := (xl + xr) / 2
	my := (yl + yr) / 2
	if node.c[0] == nil {
		node.c[0] = &Node{}
		node.c[1] = &Node{}
		node.c[2] = &Node{}
		node.c[3] = &Node{}
	}
	node.c[0].update(x1, x2, y1, y2, xl, mx, yl, my)
	node.c[1].update(x1, x2, y1, y2, xl, mx, my, yr)
	node.c[2].update(x1, x2, y1, y2, mx, xr, yl, my)
	node.c[3].update(x1, x2, y1, y2, mx, xr, my, yr)
	if node.c[0].covered && node.c[1].covered && node.c[2].covered && node.c[3].covered {
		node.covered = true
		node.c = [4]*Node{}
	}
}

func (node *Node) query(x1, x2, y1, y2, xl, xr, yl, yr int) bool {
	if x2 <= xl || xr <= x1 || y2 <= yl || yr <= y1 {
		return true
	}
	if node.covered {
		return true
	}
	if xl+1 == xr && yl+1 == yr {
		return node.covered
	}
	mx := (xl + xr) / 2
	my := (yl + yr) / 2
	if node.c[0] == nil {
		return false
	}
	return node.c[0].query(x1, x2, y1, y2, xl, mx, yl, my) &&
		node.c[1].query(x1, x2, y1, y2, xl, mx, my, yr) &&
		node.c[2].query(x1, x2, y1, y2, mx, xr, yl, my) &&
		node.c[3].query(x1, x2, y1, y2, mx, xr, my, yr)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	type Rect struct{ x1, y1, x2, y2 int }
	rects := make([]Rect, n)
	xs := make([]int, 0, 2*n)
	ys := make([]int, 0, 2*n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &rects[i].x1, &rects[i].y1, &rects[i].x2, &rects[i].y2)
		xs = append(xs, rects[i].x1, rects[i].x2)
		ys = append(ys, rects[i].y1, rects[i].y2)
	}
	sort.Ints(xs)
	xs = unique(xs)
	sort.Ints(ys)
	ys = unique(ys)

	toIndex := func(arr []int, v int) int {
		return sort.SearchInts(arr, v)
	}

	root := &Node{}
	visible := 1 // color 0 is always visible
	seen := make([]bool, n)
	for i := n - 1; i >= 0; i-- {
		r := rects[i]
		x1 := toIndex(xs, r.x1)
		x2 := toIndex(xs, r.x2)
		y1 := toIndex(ys, r.y1)
		y2 := toIndex(ys, r.y2)
		if !root.query(x1, x2, y1, y2, 0, len(xs)-1, 0, len(ys)-1) {
			visible++
			seen[i] = true
		}
		root.update(x1, x2, y1, y2, 0, len(xs)-1, 0, len(ys)-1)
	}
	fmt.Println(visible)
}

func unique(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}
