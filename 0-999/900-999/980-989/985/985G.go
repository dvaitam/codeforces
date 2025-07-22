package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct{ u, v int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	var A, B, C uint64
	fmt.Fscan(in, &A, &B, &C)

	adj := make([][]int, n)
	edges := make([]Edge, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		if u > v {
			u, v = v, u
		}
		edges[i] = Edge{u, v}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	prefix := make([]uint64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + uint64(i)
	}

	// total sum over all triples without restrictions
	var total uint64
	for i := 0; i < n; i++ {
		cnt := n - 1 - i
		if cnt >= 2 {
			choose2 := uint64(cnt * (cnt - 1) / 2)
			total += uint64(i) * choose2 * A
		}
	}
	for j := 0; j < n; j++ {
		left := j
		right := n - 1 - j
		if left > 0 && right > 0 {
			count := uint64(left * right)
			total += uint64(j) * count * B
		}
	}
	for k := 0; k < n; k++ {
		if k >= 2 {
			choose2 := uint64(k * (k - 1) / 2)
			total += uint64(k) * choose2 * C
		}
	}

	// contributions for edges
	var edgeSum uint64
	for _, e := range edges {
		u, v := e.u, e.v // ensure u<v
		// w<u
		cnt1 := uint64(u)
		sum1 := prefix[u]
		edgeSum += uint64(A)*sum1 + (uint64(B)*uint64(u)+uint64(C)*uint64(v))*cnt1
		// u<w<v
		cnt2 := uint64(v - u - 1)
		if cnt2 > 0 {
			sumMid := prefix[v] - prefix[u+1]
			edgeSum += uint64(A)*uint64(u)*cnt2 + uint64(B)*sumMid + uint64(C)*uint64(v)*cnt2
		}
		// v<w
		cnt3 := uint64(n - 1 - v)
		if cnt3 > 0 {
			sum3 := prefix[n] - prefix[v+1]
			edgeSum += (uint64(A)*uint64(u)+uint64(B)*uint64(v))*cnt3 + uint64(C)*sum3
		}
	}

	// pair of edges contributions
	var pairSum uint64
	for u := 0; u < n; u++ {
		less := make([]int, 0, len(adj[u]))
		greater := make([]int, 0, len(adj[u]))
		for _, v := range adj[u] {
			if v < u {
				less = append(less, v)
			} else if v > u {
				greater = append(greater, v)
			}
		}
		if len(less) > 1 {
			sort.Ints(less)
		}
		if len(greater) > 1 {
			sort.Ints(greater)
		}
		// case u smallest
		var prefixSum uint64
		for i, w := range greater {
			cnt := uint64(i)
			if cnt > 0 {
				pairSum += uint64(A)*uint64(u)*cnt + uint64(B)*prefixSum + uint64(C)*uint64(w)*cnt
			}
			prefixSum += uint64(w)
		}
		// case u largest
		prefixSum = 0
		for i, w := range less {
			cnt := uint64(i)
			if cnt > 0 {
				pairSum += uint64(A)*prefixSum + uint64(B)*uint64(w)*cnt + uint64(C)*uint64(u)*cnt
			}
			prefixSum += uint64(w)
		}
		// case u middle
		if len(less) > 0 && len(greater) > 0 {
			sumL := uint64(0)
			for _, x := range less {
				sumL += uint64(x)
			}
			sumR := uint64(0)
			for _, x := range greater {
				sumR += uint64(x)
			}
			lcnt := uint64(len(less))
			rcnt := uint64(len(greater))
			pairSum += uint64(A)*sumL*rcnt + uint64(B)*uint64(u)*lcnt*rcnt + uint64(C)*sumR*lcnt
		}
	}

	// triangle contributions
	deg := make([]int, n)
	for _, e := range edges {
		deg[e.u]++
		deg[e.v]++
	}
	outAdj := make([][]int, n)
	for _, e := range edges {
		u, v := e.u, e.v
		if deg[u] < deg[v] || (deg[u] == deg[v] && u < v) {
			outAdj[u] = append(outAdj[u], v)
		} else {
			outAdj[v] = append(outAdj[v], u)
		}
	}

	visited := make([]bool, n)
	var triSum uint64
	for u := 0; u < n; u++ {
		for _, v := range outAdj[u] {
			visited[v] = true
		}
		for _, v := range outAdj[u] {
			for _, w := range outAdj[v] {
				if visited[w] {
					a, b, c := u, v, w
					if a > b {
						a, b = b, a
					}
					if b > c {
						b, c = c, b
					}
					if a > b {
						a, b = b, a
					}
					triSum += uint64(A)*uint64(a) + uint64(B)*uint64(b) + uint64(C)*uint64(c)
				}
			}
		}
		for _, v := range outAdj[u] {
			visited[v] = false
		}
	}

	ans := total - edgeSum + pairSum - triSum
	fmt.Fprintln(out, ans)
}
