package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	maxV = 105
	maxE = 10000
	inf  = 0x3f3f3f3f
)

var (
	S, T, e int
	fir     [maxV]int
	to, nxt [maxE]int
	w       [maxE]int
	cost    [maxE]int
	dis     [maxV]int64
	q       [maxE]int
	vis     [maxV]bool
	cur     [maxV]int
	fdem    [maxV]int
	odd     []bool
	ansEdge []int
	sumNode []int
	anss    int64
)

func adde(x, y, z, cst int) {
	e++
	to[e] = y
	nxt[e] = fir[x]
	fir[x] = e
	w[e] = z
	cost[e] = cst
	e++
	to[e] = x
	nxt[e] = fir[y]
	fir[y] = e
	w[e] = 0
	cost[e] = -cst
}

func spfa() bool {
	// initialize distances
	for i := 1; i <= T; i++ {
		dis[i] = 1 << 60
		vis[i] = false
	}
	head, tail := 0, 0
	q[tail] = T
	dis[T] = 0
	vis[T] = true
	tail++
	for head < tail {
		u := q[head]
		head++
		vis[u] = false
		for i := fir[u]; i > 0; i = nxt[i] {
			// reverse edge index is i^1
			rev := i ^ 1
			v := to[i]
			if w[rev] == 0 {
				continue
			}
			nd := dis[u] + int64(cost[rev])
			if dis[v] > nd {
				dis[v] = nd
				if !vis[v] {
					vis[v] = true
					q[tail] = v
					tail++
				}
			}
		}
	}
	return dis[S] < 0
}

func dfs(u, flow int) int {
	if u == T || flow == 0 {
		return flow
	}
	vis[u] = true
	var used = flow
	for i := cur[u]; i > 0; i = nxt[i] {
		cur[u] = i
		v := to[i]
		if vis[v] || dis[v]+int64(cost[i]) != dis[u] || w[i] == 0 {
			continue
		}
		// augment
		can := flow
		if w[i] < can {
			can = w[i]
		}
		if can > used {
			can = used
		}
		f := dfs(v, can)
		if f > 0 {
			w[i] -= f
			w[i^1] += f
			used -= f
			if used == 0 {
				break
			}
		}
	}
	return flow - used
}

func MCMF() int64 {
	var flow int
	var res int64
	for spfa() {
		for i := 1; i <= T; i++ {
			cur[i] = fir[i]
			vis[i] = false
		}
		f := dfs(S, inf)
		flow += f
		res += dis[S] * int64(f)
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	S, T = 1, n
	e = 1
	odd = make([]bool, m+1)
	ansEdge = make([]int, m+1)
	sumNode = make([]int, n+2)
	for i := 1; i <= m; i++ {
		var x, y, z, t int
		fmt.Fscan(in, &x, &y, &z, &t)
		adde(x, y, z>>1, t<<1)
		if z&1 == 1 {
			fdem[x]--
			fdem[y]++
			odd[i] = true
			anss += int64(t)
		}
	}
	// handle demands
	for i := 2; i < n; i++ {
		if fdem[i]&1 != 0 {
			fmt.Fprintln(out, "Impossible")
			return
		}
		if fdem[i] > 0 {
			adde(S, i, fdem[i]>>1, -inf)
		} else if fdem[i] < 0 {
			adde(i, T, (-fdem[i])>>1, -inf)
		}
	}
	anss += MCMF()
	// gather results
	for i := 1; i <= m; i++ {
		// forward edge index is 2*i
		flowUsed := w[2*i+1] // reverse edge cap = flow
		ans := flowUsed << 1
		if odd[i] {
			ans |= 1
		}
		ansEdge[i] = ans
		// original adde stores forward at even index 2*i => to[2*i]
		u := to[2*i]
		v := to[2*i+1]
		sumNode[u] += ans
		sumNode[v] -= ans
	}
	for i := 2; i < n; i++ {
		if sumNode[i] != 0 {
			fmt.Fprintln(out, "Impossible")
			return
		}
	}
	fmt.Fprintln(out, "Possible")
	for i := 1; i <= m; i++ {
		fmt.Fprint(out, ansEdge[i], " ")
	}
	fmt.Fprintln(out)
}
