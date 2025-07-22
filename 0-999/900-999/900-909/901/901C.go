package main

import (
	"bufio"
	"fmt"
	"os"
)

type interval struct{ L, R int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	fmt.Fscan(reader, &n, &m)
	g := make([][]int, n+1)
	deg := make([]int, n+1)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		g[a] = append(g[a], b)
		g[b] = append(g[b], a)
		deg[a]++
		deg[b]++
	}
	removed := make([]bool, n+1)
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if deg[i] <= 1 {
			q = append(q, i)
			removed[i] = true
		}
	}
	for front := 0; front < len(q); front++ {
		v := q[front]
		for _, u := range g[v] {
			if removed[u] {
				continue
			}
			deg[u]--
			if deg[u] == 1 {
				removed[u] = true
				q = append(q, u)
			}
		}
	}
	visited := make([]bool, n+1)
	intervals := make([]interval, 0)
	stack := make([]int, 0)
	for i := 1; i <= n; i++ {
		if !removed[i] && !visited[i] {
			minv, maxv := i, i
			stack = append(stack, i)
			visited[i] = true
			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if v < minv {
					minv = v
				}
				if v > maxv {
					maxv = v
				}
				for _, u := range g[v] {
					if removed[u] || visited[u] {
						continue
					}
					visited[u] = true
					stack = append(stack, u)
				}
			}
			intervals = append(intervals, interval{minv, maxv})
		}
	}
	r := make([]int, n+2)
	for i := 1; i <= n; i++ {
		r[i] = n
	}
	for _, iv := range intervals {
		if iv.R-1 < r[iv.L] {
			r[iv.L] = iv.R - 1
		}
	}
	for i := n - 1; i >= 1; i-- {
		if r[i] > r[i+1] {
			r[i] = r[i+1]
		}
	}
	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pref[i] = pref[i-1] + int64(r[i]-i+1)
	}
	var qnum int
	fmt.Fscan(reader, &qnum)
	writer := bufio.NewWriter(os.Stdout)
	for ; qnum > 0; qnum-- {
		var L, R int
		fmt.Fscan(reader, &L, &R)
		low, high := L, R
		p := L - 1
		for low <= high {
			mid := (low + high) / 2
			if r[mid] <= R {
				p = mid
				low = mid + 1
			} else {
				high = mid - 1
			}
		}
		ans := int64(0)
		if p >= L {
			ans += pref[p] - pref[L-1]
		}
		cnt := int64(R - p)
		ans += cnt * (cnt + 1) / 2
		fmt.Fprintln(writer, ans)
	}
	writer.Flush()
}
