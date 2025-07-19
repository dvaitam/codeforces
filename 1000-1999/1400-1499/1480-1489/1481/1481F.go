package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

type Item struct {
	val, cnt int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	line, _ := in.ReadString('\n')
	parts := strings.Fields(line)
	n, _ := strconv.Atoi(parts[0])
	k, _ := strconv.Atoi(parts[1])
	par := make([]int, n)
	adj := make([][]int, n)
	for i := 1; i < n; i++ {
		line, _ = in.ReadString('\n')
		p, _ := strconv.Atoi(strings.TrimSpace(line))
		p--
		par[i] = p
		adj[p] = append(adj[p], i)
	}
	dep := make([]int, n)
	v := make([][]int, n)
	// BFS
	q := []int{0}
	for idx := 0; idx < len(q); idx++ {
		u := q[idx]
		v[dep[u]] = append(v[dep[u]], u)
		for _, to := range adj[u] {
			dep[to] = dep[u] + 1
			q = append(q, to)
		}
	}
	maxd := 0
	for i := 0; i < n; i++ {
		if dep[i] > maxd {
			maxd = dep[i]
		}
	}
	// group depths by size
	deps := make(map[int][]int)
	cnt := make(map[int]int)
	for d := 0; d <= maxd; d++ {
		sz := len(v[d])
		if sz > 0 {
			deps[sz] = append(deps[sz], d)
			cnt[sz]++
		}
	}
	if k == 0 || k == n {
		// simple
		val := maxd + 1
		if k != 0 {
			fmt.Fprintln(out, val)
			fmt.Fprintln(out, strings.Repeat("a", n))
		} else {
			fmt.Fprintln(out, val)
			fmt.Fprintln(out, strings.Repeat("b", n))
		}
		return
	}
	// build items for knapsack
	items := make([]Item, 0)
	for size, c := range cnt {
		m := c
		p := 1
		for m > 0 {
			take := p
			if take > m {
				take = m
			}
			items = append(items, Item{val: size, cnt: take})
			m -= take
			p <<= 1
		}
	}
	// dp bitset
	L := (k + 64) >> 6
	vis := make([]uint64, L)
	vis[0] = 1
	prv := make([]Item, k+1)
	// dp
	for _, it := range items {
		w := it.val * it.cnt
		// shift vis by w into nxt
		nxt := make([]uint64, L)
		copy(nxt, vis)
		block := w >> 6
		shift := uint(w & 63)
		for j := L - 1; j >= 0; j-- {
			src := j - block
			if src < 0 {
				continue
			}
			vsrc := vis[src]
			var sv uint64
			if shift == 0 {
				sv = vsrc
			} else {
				sv = vsrc << shift
				if src > 0 {
					sv |= vis[src-1] >> (64 - shift)
				}
			}
			nxt[j] |= sv
		}
		// record new bits
		for j := 0; j < L; j++ {
			newb := nxt[j] ^ vis[j]
			for newb != 0 {
				t := bits.TrailingZeros64(newb)
				idx := j<<6 | t
				if idx <= k {
					prv[idx] = it
				}
				newb &= newb - 1
			}
		}
		vis = nxt
	}
	ans := make([]bool, n)
	if (vis[k>>6]>>(uint(k)&63))&1 == 1 {
		// can knapsack
		fmt.Fprintln(out, maxd+1)
		cur := k
		for cur > 0 {
			it := prv[cur]
			for t := 0; t < it.cnt; t++ {
				size := it.val
				ds := deps[size]
				d := ds[len(ds)-1]
				deps[size] = ds[:len(ds)-1]
				for _, x := range v[d] {
					ans[x] = true
				}
			}
			cur -= it.val * it.cnt
		}
	} else {
		// greedy
		fmt.Fprintln(out, maxd+2)
		m := n
		rem := k
		for d := 0; m > 0; d++ {
			sz := len(v[d])
			if sz <= rem {
				rem -= sz
				for _, x := range v[d] {
					ans[x] = true
				}
			} else if sz > m-rem {
				// need split
				p := rem < m-rem
				c := rem
				if m-rem < rem {
					c = m - rem
				}
				if p {
					rem -= c
				}
				leaves := make([]int, 0, sz)
				for _, x := range v[d] {
					if len(adj[x]) == 0 {
						leaves = append(leaves, x)
					}
				}
				// assign p to some leaves
				for i := 0; i < c; i++ {
					x := leaves[len(leaves)-1]
					leaves = leaves[:len(leaves)-1]
					ans[x] = p
				}
				// rest leaves and non-leaves
				for _, x := range leaves {
					ans[x] = !p
				}
				for _, x := range v[d] {
					if len(adj[x]) > 0 {
						ans[x] = !p
					}
				}
			}
			m -= sz
		}
	}
	// output chars
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if ans[i] {
			b[i] = 'a'
		} else {
			b[i] = 'b'
		}
	}
	fmt.Fprintln(out, string(b))
}
