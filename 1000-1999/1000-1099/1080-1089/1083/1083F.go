package main

import (
	"bufio"
	"fmt"
	"os"
)

const B = 500

var n, k, q, m, zch int
var a, b, d []int
var pos, epos, spos []int
var off []int
var cnt []map[int]int

func modify(p, v int) {
	if p < 0 || p > n {
		return
	}
	l := pos[p]
	r := epos[p]
	// remove old end zero status
	if r > 0 {
		br := (r - 1) / B
		if spos[r-1]^off[br] == 0 {
			zch--
		}
	}
	idl := l / B
	idr := r / B
	if idl == idr {
		for i := l; i < r; i++ {
			blk := idl
			cnt[blk][spos[i]]--
			spos[i] ^= v
			cnt[blk][spos[i]]++
		}
	} else {
		endL := (idl + 1) * B
		for i := l; i < endL; i++ {
			blk := idl
			cnt[blk][spos[i]]--
			spos[i] ^= v
			cnt[blk][spos[i]]++
		}
		for blk := idl + 1; blk < idr; blk++ {
			off[blk] ^= v
		}
		startR := idr * B
		for i := startR; i < r; i++ {
			blk := idr
			cnt[blk][spos[i]]--
			spos[i] ^= v
			cnt[blk][spos[i]]++
		}
	}
	// add new end zero status
	if r > 0 {
		br := (r - 1) / B
		if spos[r-1]^off[br] == 0 {
			zch++
		}
	}
}

func output(w *bufio.Writer) {
	if zch != k {
		w.WriteString("-1\n")
	} else {
		sumZero := 0
		for blk := range off {
			sumZero += cnt[blk][off[blk]]
		}
		// total pairs = m, answer = m - zeroes
		res := m - sumZero
		w.WriteString(fmt.Sprintf("%d\n", res))
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &k, &q)
	a = make([]int, n+2)
	b = make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	d = make([]int, n+2)
	nch := make([][]int, k)
	for i := 0; i < k; i++ {
		nch[i] = make([]int, 0)
	}
	for i := 0; i <= n; i++ {
		d[i] = (a[i] ^ a[i+1]) ^ (b[i] ^ b[i+1])
		rem := i % k
		nch[rem] = append(nch[rem], i)
	}
	pos = make([]int, n+2)
	epos = make([]int, n+2)
	spos = make([]int, n+2)
	m = 0
	zch = 0
	for i := 0; i < k; i++ {
		s := 0
		for _, u := range nch[i] {
			s ^= d[u]
			spos[m] = s
			pos[u] = m
			m++
		}
		if m > 0 && spos[m-1] == 0 {
			zch++
		}
		for _, u := range nch[i] {
			epos[u] = m
		}
	}
	blocks := (m + B - 1) / B
	off = make([]int, blocks)
	cnt = make([]map[int]int, blocks)
	for i := 0; i < blocks; i++ {
		cnt[i] = make(map[int]int)
	}
	for i := 0; i < m; i++ {
		blk := i / B
		cnt[blk][spos[i]]++
	}
	output(writer)
	for qi := 0; qi < q; qi++ {
		var op string
		var p, v int
		fmt.Fscan(reader, &op, &p, &v)
		if op[0] == 'a' {
			a[p] ^= v
			v, a[p] = a[p], v
		} else {
			b[p] ^= v
			v, b[p] = b[p], v
		}
		modify(p, v)
		modify(p-1, v)
		output(writer)
	}
}
