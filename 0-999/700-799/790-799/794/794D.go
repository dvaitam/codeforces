package main

import (
	"bufio"
	"fmt"
	"os"
)

var n, m int
var ghead, gnxt, gto, deg []int
var done, dead []bool
var prv, q, path, id, cnt []int

func solve() bool {
	s, t := -1, -1
	// find non-edge s->t
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			done[j] = false
		}
		for x := ghead[i]; x != -1; x = gnxt[x] {
			done[gto[x]] = true
		}
		for j := 0; j < n; j++ {
			if i != j && !done[j] {
				s, t = i, j
				break
			}
		}
		if s != -1 {
			break
		}
	}
	if s == -1 {
		for i := 0; i < n; i++ {
			id[i] = 0
		}
		return true
	}
	// BFS from s
	for i := 0; i < n; i++ {
		done[i] = false
		prv[i] = -2
	}
	qhead, qtail := 0, 0
	done[s] = true
	prv[s] = -1
	q[qhead] = s
	qhead++
	for qtail < qhead {
		at := q[qtail]
		qtail++
		for x := ghead[at]; x != -1; x = gnxt[x] {
			to := gto[x]
			if !done[to] {
				done[to] = true
				prv[to] = at
				q[qhead] = to
				qhead++
			}
		}
	}
	// mark path nodes
	for i := 0; i < n; i++ {
		dead[i] = false
		done[i] = false
	}
	for at := t; at != -1; at = prv[at] {
		done[at] = true
	}
	// mark dead from internal path nodes
	for at := prv[t]; at != s; at = prv[at] {
		for x := ghead[at]; x != -1; x = gnxt[x] {
			dead[gto[x]] = true
		}
	}
	// extend from s
	for {
		to := -1
		for x := ghead[s]; x != -1; x = gnxt[x] {
			v := gto[x]
			if !done[v] && !dead[v] {
				to = v
				break
			}
		}
		if to == -1 {
			break
		}
		for x := ghead[s]; x != -1; x = gnxt[x] {
			dead[gto[x]] = true
		}
		done[to] = true
		prv[s] = to
		prv[to] = -1
		s = to
	}
	// extend from t
	for {
		to := -1
		for x := ghead[t]; x != -1; x = gnxt[x] {
			v := gto[x]
			if !done[v] && !dead[v] {
				to = v
				break
			}
		}
		if to == -1 {
			break
		}
		for x := ghead[t]; x != -1; x = gnxt[x] {
			dead[gto[x]] = true
		}
		done[to] = true
		prv[to] = t
		t = to
	}
	// build path
	npath := 0
	for at := t; at != -1; at = prv[at] {
		path[npath] = at
		npath++
	}
	// reverse path
	for i, j := 0, npath-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	// assign ids on path
	for i := 0; i < n; i++ {
		id[i] = -1
	}
	for i := 0; i < npath; i++ {
		id[path[i]] = i
	}
	// assign other nodes
	for at := 0; at < n; at++ {
		if id[at] != -1 {
			continue
		}
		mn, mx := n, -1
		any := false
		for x := ghead[at]; x != -1; x = gnxt[x] {
			v := gto[x]
			if id[v] != -1 {
				if !any {
					mn, mx = id[v], id[v]
					any = true
				} else {
					if id[v] < mn {
						mn = id[v]
					}
					if id[v] > mx {
						mx = id[v]
					}
				}
			}
		}
		if !any {
			return false
		}
		switch mx - mn {
		case 2:
			id[at] = mn + 1
		case 1:
			if mn == 0 {
				id[at] = mn
			} else if mx == npath-1 {
				id[at] = mx
			} else {
				return false
			}
		default:
			return false
		}
	}
	// check edges
	for u := 0; u < n; u++ {
		for x := ghead[u]; x != -1; x = gnxt[x] {
			v := gto[x]
			if id[u]-id[v] > 1 || id[v]-id[u] > 1 {
				return false
			}
		}
	}
	// count segments
	for i := 0; i < npath; i++ {
		cnt[i] = 0
	}
	for u := 0; u < n; u++ {
		cnt[id[u]]++
	}
	// degree check
	for u := 0; u < n; u++ {
		want := -1
		for d := -1; d <= 1; d++ {
			nid := id[u] + d
			if nid >= 0 && nid < npath {
				want += cnt[nid]
			}
		}
		if want != deg[u] {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fscan(in, &n, &m)
	ghead = make([]int, n)
	gnxt = make([]int, 2*m)
	gto = make([]int, 2*m)
	deg = make([]int, n)
	done = make([]bool, n)
	dead = make([]bool, n)
	prv = make([]int, n)
	q = make([]int, n)
	path = make([]int, n)
	id = make([]int, n)
	cnt = make([]int, n)
	for i := 0; i < n; i++ {
		ghead[i] = -1
	}
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a--
		b--
		gnxt[2*i] = ghead[a]
		ghead[a] = 2 * i
		gto[2*i] = b
		deg[a]++
		gnxt[2*i+1] = ghead[b]
		ghead[b] = 2*i + 1
		gto[2*i+1] = a
		deg[b]++
	}
	if !solve() {
		fmt.Fprintln(out, "NO")
		return
	}
	fmt.Fprintln(out, "YES")
	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, id[i]+1)
	}
	out.WriteByte('\n')
}
