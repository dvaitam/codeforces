package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

type edge struct {
	u, v int
	w    int64
}

type dsu struct {
	p []int
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	for i := range p {
		p[i] = i
	}
	return &dsu{p: p}
}

func (d *dsu) find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.find(d.p[x])
	}
	return d.p[x]
}

const infAns int64 = 1 << 62

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m, p int
		fmt.Fscan(in, &n, &m, &p)

		required := make([]bool, 2*n+5)
		for i := 0; i < p; i++ {
			var x int
			fmt.Fscan(in, &x)
			required[x] = true
		}

		edges := make([]edge, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &edges[i].u, &edges[i].v, &edges[i].w)
		}
		sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })

		maxNodes := 2*n + 5
		parent := newDSU(maxNodes)
		id := make([]int, maxNodes)
		for i := 1; i <= n; i++ {
			id[i] = i
		}
		children := make([][2]int, maxNodes)
		weight := make([]int64, maxNodes)
		cur := n
		for _, e := range edges {
			ru, rv := parent.find(e.u), parent.find(e.v)
			if ru == rv {
				continue
			}
			cur++
			children[cur] = [2]int{id[ru], id[rv]}
			weight[cur] = e.w
			parent.p[ru] = cur
			parent.p[rv] = cur
			parent.p[cur] = cur
			id[cur] = cur
		}
		root := id[parent.find(1)]

		// Build postorder for fast per-lambda evaluation
		order := make([]int, 0, cur)
		stack := []int{root}
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			order = append(order, v)
			if v > n {
				stack = append(stack, children[v][0], children[v][1])
			}
		}
		// reverse for postorder
		for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
			order[i], order[j] = order[j], order[i]
		}

		cnt := make([]int, cur+1)
		for _, v := range order {
			if v <= n {
				if required[v] {
					cnt[v] = 1
				}
			} else {
				l, r := children[v][0], children[v][1]
				cnt[v] = cnt[l] + cnt[r]
			}
		}

		type res struct {
			val float64
			k   int
		}

		costArr := make([]float64, cur+1)
		srvArr := make([]int, cur+1)

		evaluate := func(lambda float64) res {
			for _, v := range order {
				if v <= n {
					if required[v] {
						costArr[v] = lambda
						srvArr[v] = 1
					} else {
						costArr[v] = 0
						srvArr[v] = 0
					}
					continue
				}
				l, r := children[v][0], children[v][1]
				cl, sl := costArr[l], srvArr[l]
				cr, sr := costArr[r], srvArr[r]
				bestVal := cl + cr
				bestK := sl + sr
				if cnt[r] > 0 {
					val := cl + float64(cnt[r])*float64(weight[v])
					if val < bestVal || (math.Abs(val-bestVal) < 1e-12 && sl < bestK) {
						bestVal = val
						bestK = sl
					}
				}
				if cnt[l] > 0 {
					val := cr + float64(cnt[l])*float64(weight[v])
					if val < bestVal || (math.Abs(val-bestVal) < 1e-12 && sr < bestK) {
						bestVal = val
						bestK = sr
					}
				}
				costArr[v] = bestVal
				srvArr[v] = bestK
			}
			return res{val: costArr[root], k: srvArr[root]}
		}

		ans := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			ans[i] = infAns
		}

		addAns := func(k int, val float64) {
			if k <= 0 || k > n {
				return
			}
			iv := int64(math.Round(val))
			if iv < ans[k] {
				ans[k] = iv
			}
		}

		var dfs func(float64, res, float64, res)
		dfs = func(lamL float64, left res, lamR float64, right res) {
			addAns(left.k, left.val-lamL*float64(left.k))
			addAns(right.k, right.val-lamR*float64(right.k))
			if left.k == right.k {
				return
			}
			mid := (left.val - right.val) / float64(right.k-left.k)
			midRes := evaluate(mid)
			addAns(midRes.k, midRes.val-mid*float64(midRes.k))
			if midRes.k == left.k || midRes.k == right.k {
				return
			}
			dfs(lamL, left, mid, midRes)
			dfs(mid, midRes, lamR, right)
		}

		lRes := evaluate(-1e12)
		rRes := evaluate(1e12)
		dfs(-1e12, lRes, 1e12, rRes)

		for i := p; i <= n; i++ {
			ans[i] = 0
		}
		for i := n - 1; i >= 1; i-- {
			if ans[i] > ans[i+1] {
				ans[i] = ans[i+1]
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
