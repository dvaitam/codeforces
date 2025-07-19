package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxN = 40000

type pair struct{ first, second int }

var (
	n, k   int
	d      [][]int
	who    []int
	parent []int
	S      map[int64]struct{}
	v      [][]int
	nd     [][]int
	dep    []int
	gt     []bool
	wt     [][]pair
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func imp() {
	fmt.Fprintln(writer, -1)
	writer.Flush()
	os.Exit(0)
}

func find(x int) int {
	if parent[x] != x {
		parent[x] = find(parent[x])
	}
	return parent[x]
}

func union(x, y int) {
	fx := find(x)
	fy := find(y)
	parent[fx] = fy
}

func addEdge(a, b int) {
	if a > b {
		a, b = b, a
	}
	key := int64(a)*maxN + int64(b)
	if _, exists := S[key]; exists {
		return
	}
	if find(a) == find(b) {
		imp()
	}
	gt[a], gt[b] = true, true
	union(a, b)
	S[key] = struct{}{}
	v[a] = append(v[a], b)
	v[b] = append(v[b], a)
}

func goDFS(now, prt, tdep int) {
	if dep[now] != 0 {
		imp()
	}
	dep[now] = tdep
	for _, son := range v[now] {
		if son == prt {
			continue
		}
		goDFS(son, now, tdep+1)
	}
}

func gogo(id, now, prt, tdep int) {
	if d[id][now] != tdep {
		imp()
	}
	for _, son := range v[now] {
		if son == prt {
			continue
		}
		gogo(id, son, now, tdep+1)
	}
}

func cons(id int) {
	if len(wt[id]) == 0 {
		return
	}
	sort.Slice(wt[id], func(i, j int) bool {
		if wt[id][i].first != wt[id][j].first {
			return wt[id][i].first < wt[id][j].first
		}
		return wt[id][i].second < wt[id][j].second
	})
	// collect unique first values
	dd := make([]int, len(wt[id]))
	for i, p := range wt[id] {
		dd[i] = p.first
	}
	m := 0
	for i := range dd {
		if m == 0 || dd[i] != dd[m-1] {
			dd[m] = dd[i]
			m++
		}
	}
	if m == 0 {
		imp()
	}
	for i := 0; i < m; i++ {
		if dd[i] != i+1 {
			imp()
		}
	}
	pre := id
	for i := range wt[id] {
		addEdge(pre, wt[id][i].second)
		if i+1 < len(wt[id]) && wt[id][i].first != wt[id][i+1].first {
			pre = wt[id][i].second
		}
	}
}

func recheck() {
	if len(S)+1 != n {
		imp()
	}
	for i := 0; i < k; i++ {
		gogo(i, who[i], who[i], 0)
	}
}

func output() {
	for key := range S {
		u := int(key / maxN)
		v2 := int(key % maxN)
		fmt.Fprintln(writer, u, v2)
	}
	writer.Flush()
	os.Exit(0)
}

func solve2() {
	for i := 1; i < k; i++ {
		if d[0][who[i]] != d[i][who[0]] {
			imp()
		}
		dst := d[0][who[i]]
		for j := 1; j <= n; j++ {
			if d[0][j]+d[i][j] == dst {
				if nd[i][d[0][j]] != 0 {
					imp()
				}
				nd[i][d[0][j]] = j
			}
		}
		for j := 1; j <= dst; j++ {
			if nd[i][j] == 0 {
				imp()
			}
			addEdge(nd[i][j-1], nd[i][j])
		}
	}
	goDFS(who[0], -1, 1)
	for i := 1; i <= n; i++ {
		if gt[i] {
			continue
		}
		bst := who[0]
		for j := 1; j < k; j++ {
			tmp := who[0]
			dlt := abs(d[0][i] - d[j][i])
			if dlt == d[0][who[j]] {
				if d[0][i] > d[j][i] {
					tmp = who[j]
				}
			} else {
				if dlt%2 != d[0][who[j]]%2 {
					imp()
				}
				dlt = d[j][i] - d[0][i]
				at := (d[0][who[j]] - dlt) / 2
				tmp = nd[j][at]
			}
			if dep[tmp] > dep[bst] {
				bst = tmp
			}
		}
		wt[bst] = append(wt[bst], pair{d[0][i] - dep[bst] + 1, i})
	}
	for i := 1; i <= n; i++ {
		cons(i)
	}
	recheck()
	output()
}

func solve() {
	if k == 1 {
		vv := make([]int, n)
		for i := 1; i <= n; i++ {
			vv[i-1] = d[0][i]
		}
		sort.Ints(vv)
		if vv[0] != 0 {
			imp()
		}
		if n > 1 && vv[1] == 0 {
			imp()
		}
		m := 0
		for i := range vv {
			if m == 0 || vv[i] != vv[m-1] {
				vv[m] = vv[i]
				m++
			}
		}
		for i := 0; i < m; i++ {
			if vv[i] != i {
				imp()
			}
		}
		v2 := make([]pair, n)
		for i := 1; i <= n; i++ {
			v2[i-1] = pair{d[0][i], i}
		}
		sort.Slice(v2, func(i, j int) bool {
			if v2[i].first != v2[j].first {
				return v2[i].first < v2[j].first
			}
			return v2[i].second < v2[j].second
		})
		pre := v2[0].second
		for i := 1; i < len(v2); i++ {
			addEdge(pre, v2[i].second)
			if i+1 < len(v2) && v2[i].first != v2[i+1].first {
				pre = v2[i].second
			}
		}
		output()
	}
	solve2()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	defer writer.Flush()
	fmt.Fscan(reader, &n, &k)
	d = make([][]int, k)
	who = make([]int, k)
	for i := 0; i < k; i++ {
		d[i] = make([]int, n+1)
	}
	nd = make([][]int, k)
	for i := 0; i < k; i++ {
		nd[i] = make([]int, n+1)
	}
	parent = make([]int, n+1)
	S = make(map[int64]struct{})
	v = make([][]int, n+1)
	dep = make([]int, n+1)
	gt = make([]bool, n+1)
	wt = make([][]pair, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
	}
	for i := 0; i < k; i++ {
		for j := 1; j <= n; j++ {
			fmt.Fscan(reader, &d[i][j])
			if d[i][j] == 0 {
				if who[i] != 0 {
					imp()
				}
				who[i] = j
			}
		}
	}
	solve()
}
