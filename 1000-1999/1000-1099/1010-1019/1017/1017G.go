package main

import (
   "bufio"
   "os"
)

var (
	n, m                 int
	p                    [][]int
	fa, dep, sz, ws, top []int
	id, idr, ir, bv      []int
	idm                  int
	f                    []node
	ru                   []pair
)

type node struct{ mi, ui, ua, tag int }
type pair struct{ x, v int }

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	read := func() int {
		var x int
		var c byte
		var neg bool
		for {
			b, err := reader.ReadByte()
			if err != nil {
				break
			}
			c = b
			if c == '-' || (c >= '0' && c <= '9') {
				break
			}
		}
		if c == '-' {
			neg = true
			c, _ = reader.ReadByte()
		}
		for ; c >= '0' && c <= '9'; c, _ = reader.ReadByte() {
			x = x*10 + int(c-'0')
		}
		if neg {
			x = -x
		}
		return x
	}
	n = read()
	m = read()
	p = make([][]int, n+1)
	fa = make([]int, n+1)
	dep = make([]int, n+1)
	sz = make([]int, n+1)
	ws = make([]int, n+1)
	top = make([]int, n+1)
	id = make([]int, n+1)
	idr = make([]int, n+1)
	ir = make([]int, n+1)
	bv = make([]int, n+2)
	for i := 2; i <= n; i++ {
		fa[i] = read()
		p[fa[i]] = append(p[fa[i]], i)
	}
	dep[1] = 0
	dfs(1)
	dfs2(1, 1)
	for i := 1; i <= n; i++ {
		ir[id[i]] = i
		dep[i] = dep[fa[i]] + 1
		bv[id[i]] = -dep[i]
	}
	f = make([]node, 4*(n+2))
	ru = make([]pair, n+5)
	build(1, 1, n)
	out := func(s string) {
		writer.WriteString(s)
	}
	for ; m > 0; m-- {
		op := read()
		x := read()
		if op == 1 {
			modu(1, 1, n, id[x], 1)
			mod(1, 1, n, id[x]+1, idr[x], 1)
		} else if op == 2 {
			rc := 0
			rst(1, 1, n, id[x], idr[x], &rc)
			for i := 0; i < rc; i++ {
				xx := ru[i].x
				vv := ru[i].v
				mod(1, 1, n, id[xx]+1, idr[xx], -vv)
			}
			y := x
			r := int(1e9)
			for y > 0 {
				v := qry(1, 1, n, id[top[y]], id[y])
				if v < r {
					r = v
				}
				y = fa[top[y]]
			}
			r -= qry(1, 1, n, id[x], id[x])
			modu(1, 1, n, id[x], r)
			mod(1, 1, n, id[x]+1, idr[x], r)
		} else {
			v := qry(1, 1, n, id[x], id[x]) + qryu(1, 1, n, id[x])
			r := v
			y := x
			for y > 0 {
				v2 := qry(1, 1, n, id[top[y]], id[y])
				if v2 < r {
					r = v2
				}
				y = fa[top[y]]
			}
			if r < v {
				out("black\n")
			} else {
				out("white\n")
			}
		}
	}
}

// HLD and segment tree functions
func dfs(x int) {
	sz[x] = 1
	maxsz := 0
	for _, c := range p[x] {
		dep[c] = dep[x] + 1
		dfs(c)
		sz[x] += sz[c]
		if sz[c] > maxsz {
			maxsz = sz[c]
			ws[x] = c
		}
	}
}

func dfs2(x, t int) {
	top[x] = t
	idm++
	id[x] = idm
	if ws[x] != 0 {
		dfs2(ws[x], t)
	}
	for _, c := range p[x] {
		if c != ws[x] {
			dfs2(c, c)
		}
	}
	idr[x] = idm
}

func build(o, l, r int) {
	if l == r {
		f[o].mi = bv[l]
		return
	}
	m := (l + r) >> 1
	build(o<<1, l, m)
	build(o<<1|1, m+1, r)
	pu(o)
}

func st(o, v int) {
	f[o].mi += v
	f[o].tag += v
}

func pd(o int) {
	if f[o].tag != 0 {
		st(o<<1, f[o].tag)
		st(o<<1|1, f[o].tag)
		f[o].tag = 0
	}
}

func pu(o int) {
	a, b := o<<1, o<<1|1
	if f[a].mi < f[b].mi {
		f[o].mi = f[a].mi
	} else {
		f[o].mi = f[b].mi
	}
	if f[a].ui < f[b].ui {
		f[o].ui = f[a].ui
	} else {
		f[o].ui = f[b].ui
	}
	if f[a].ua > f[b].ua {
		f[o].ua = f[a].ua
	} else {
		f[o].ua = f[b].ua
	}
}

func mod(o, l, r, ql, qr, v int) {
	if ql > qr {
		return
	}
	_mod(o, l, r, ql, qr, v)
}
func _mod(o, l, r, ql, qr, v int) {
	if ql <= l && r <= qr {
		st(o, v)
		return
	}
	pd(o)
	m := (l + r) >> 1
	if ql <= m {
		_mod(o<<1, l, m, ql, qr, v)
	}
	if qr > m {
		_mod(o<<1|1, m+1, r, ql, qr, v)
	}
	pu(o)
}

func modu(o, l, r, p, v int) {
	if l == r {
		f[o].ui += v
		f[o].ua += v
		return
	}
	pd(o)
	m := (l + r) >> 1
	if p <= m {
		modu(o<<1, l, m, p, v)
	} else {
		modu(o<<1|1, m+1, r, p, v)
	}
	pu(o)
}

func qry(o, l, r, ql, qr int) int {
	res := int(1e9)
	return _qry(o, l, r, ql, qr, res)
}
func _qry(o, l, r, ql, qr, res int) int {
	if ql <= l && r <= qr {
		if f[o].mi < res {
			res = f[o].mi
		}
		return res
	}
	pd(o)
	m := (l + r) >> 1
	if ql <= m {
		res = _qry(o<<1, l, m, ql, qr, res)
	}
	if qr > m {
		res = _qry(o<<1|1, m+1, r, ql, qr, res)
	}
	return res
}

func qryu(o, l, r, p int) int {
	if l == r {
		return f[o].ua
	}
	pd(o)
	m := (l + r) >> 1
	if p <= m {
		return qryu(o<<1, l, m, p)
	}
	return qryu(o<<1|1, m+1, r, p)
}

func rst(o, l, r, ql, qr int, rc *int) {
   if ql <= l && r <= qr && f[o].ua == 0 && f[o].ui == 0 {
       return
   }
   if l == r {
       idx := *rc
       ru[idx] = pair{ir[l], f[o].ua}
       *rc = idx + 1
       f[o].ua = 0
       f[o].ui = 0
       return
   }
	pd(o)
	m := (l + r) >> 1
	if ql <= m {
		rst(o<<1, l, m, ql, qr, rc)
	}
	if qr > m {
		rst(o<<1|1, m+1, r, ql, qr, rc)
	}
	pu(o)
}
