package main

import (
	"io"
	"os"
	"strconv"
)

type Edge struct {
	to int
	w  int64
}

type Frame struct {
	u    int
	next int
}

func lowerBound(a []int64, x int64) int {
	l, r := 0, len(a)
	for l < r {
		m := (l + r) >> 1
		if a[m] >= x {
			r = m
		} else {
			l = m + 1
		}
	}
	return l
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	idx := 0
	nextInt := func() int64 {
		for idx < len(data) && (data[idx] < '0' || data[idx] > '9') {
			idx++
		}
		var v int64
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			v = v*10 + int64(data[idx]-'0')
			idx++
		}
		return v
	}

	n := int(nextInt())
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = nextInt()
	}

	parent := make([]int, n+1)
	children := make([][]Edge, n+1)
	for v := 2; v <= n; v++ {
		p := int(nextInt())
		w := nextInt()
		parent[v] = p
		children[p] = append(children[p], Edge{to: v, w: w})
	}

	dist := make([]int64, n+1)
	delta := make([]int, n+1)
	order := make([]int, 0, n)
	pathNodes := make([]int, 0, n)
	pathDist := make([]int64, 0, n)

	stack := make([]Frame, 0, n)
	stack = append(stack, Frame{u: 1, next: -1})

	for len(stack) > 0 {
		top := len(stack) - 1
		if stack[top].next == -1 {
			u := stack[top].u
			order = append(order, u)
			pathNodes = append(pathNodes, u)
			pathDist = append(pathDist, dist[u])

			pos := lowerBound(pathDist, dist[u]-a[u])
			x := pathNodes[pos]
			delta[parent[u]]++
			delta[parent[x]]--

			stack[top].next = 0
		}

		u := stack[top].u
		if stack[top].next < len(children[u]) {
			e := children[u][stack[top].next]
			stack[top].next++
			dist[e.to] = dist[u] + e.w
			stack = append(stack, Frame{u: e.to, next: -1})
		} else {
			pathNodes = pathNodes[:len(pathNodes)-1]
			pathDist = pathDist[:len(pathDist)-1]
			stack = stack[:len(stack)-1]
		}
	}

	sum := make([]int, n+1)
	ans := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		sum[u] += delta[u]
		ans[u] = sum[u]
		if parent[u] != 0 {
			sum[parent[u]] += sum[u]
		}
	}

	out := make([]byte, 0, n*7)
	for i := 1; i <= n; i++ {
		if i > 1 {
			out = append(out, ' ')
		}
		out = strconv.AppendInt(out, int64(ans[i]), 10)
	}
	out = append(out, '\n')
	os.Stdout.Write(out)
}