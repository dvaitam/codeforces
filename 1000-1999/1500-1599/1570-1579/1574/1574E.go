package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 998244353
const SHIFT = 21

type RCData struct {
	cnt      [2][2]int
	conflict bool
}

var (
	n, m, k     int
	rows        map[int]*RCData
	cols        map[int]*RCData
	assignments map[int64]int
	rowConflict int
	colConflict int
	rowActive   int
	colActive   int
	parCount    [2][2]int
)

func key(x, y int) int64 {
	return int64(x)<<SHIFT | int64(y)
}

func updateRow(x int) {
	d := rows[x]
	conflict := false
	for p := 0; p < 2; p++ {
		if d.cnt[p][0] > 0 && d.cnt[p][1] > 0 {
			conflict = true
		}
	}
	val0, val1 := -1, -1
	if d.cnt[0][0] > 0 {
		val0 = 0
	}
	if d.cnt[0][1] > 0 {
		val0 = 1
	}
	if d.cnt[1][0] > 0 {
		val1 = 0
	}
	if d.cnt[1][1] > 0 {
		val1 = 1
	}
	if val0 != -1 && val1 != -1 && (val0^val1) != 1 {
		conflict = true
	}
	if conflict && !d.conflict {
		rowConflict++
		d.conflict = true
	} else if !conflict && d.conflict {
		rowConflict--
		d.conflict = false
	}
}

func updateCol(y int) {
	d := cols[y]
	conflict := false
	for p := 0; p < 2; p++ {
		if d.cnt[p][0] > 0 && d.cnt[p][1] > 0 {
			conflict = true
		}
	}
	val0, val1 := -1, -1
	if d.cnt[0][0] > 0 {
		val0 = 0
	}
	if d.cnt[0][1] > 0 {
		val0 = 1
	}
	if d.cnt[1][0] > 0 {
		val1 = 0
	}
	if d.cnt[1][1] > 0 {
		val1 = 1
	}
	if val0 != -1 && val1 != -1 && (val0^val1) != 1 {
		conflict = true
	}
	if conflict && !d.conflict {
		colConflict++
		d.conflict = true
	} else if !conflict && d.conflict {
		colConflict--
		d.conflict = false
	}
}

func add(x, y, t int) {
	assignments[key(x, y)] = t
	rp := y % 2
	cp := x % 2
	if rows[x] == nil {
		rows[x] = &RCData{}
		rowActive++
	}
	if cols[y] == nil {
		cols[y] = &RCData{}
		colActive++
	}
	rows[x].cnt[rp][t]++
	cols[y].cnt[cp][t]++
	updateRow(x)
	updateCol(y)
	parCount[(x%2)^(y%2)][t]++
}

func removeCell(x, y int) {
	kxy := key(x, y)
	t, ok := assignments[kxy]
	if !ok {
		return
	}
	delete(assignments, kxy)
	rp := y % 2
	cp := x % 2
	rows[x].cnt[rp][t]--
	cols[y].cnt[cp][t]--
	updateRow(x)
	updateCol(y)
	parCount[(x%2)^(y%2)][t]--
	if rows[x].cnt[0][0]+rows[x].cnt[0][1]+rows[x].cnt[1][0]+rows[x].cnt[1][1] == 0 {
		if rows[x].conflict {
			rowConflict--
		}
		delete(rows, x)
		rowActive--
	}
	if cols[y].cnt[0][0]+cols[y].cnt[0][1]+cols[y].cnt[1][0]+cols[y].cnt[1][1] == 0 {
		if cols[y].conflict {
			colConflict--
		}
		delete(cols, y)
		colActive--
	}
}

func powmod(a int64, b int) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func waysRowsAlt() int64 {
	if colConflict > 0 {
		return 0
	}
	free := m - colActive
	return powmod(2, free)
}

func waysColsAlt() int64 {
	if rowConflict > 0 {
		return 0
	}
	free := n - rowActive
	return powmod(2, free)
}

func waysBoth() int64 {
	valid0 := parCount[0][1] == 0 && parCount[1][0] == 0
	valid1 := parCount[0][0] == 0 && parCount[1][1] == 0
	cnt := int64(0)
	if valid0 {
		cnt++
	}
	if valid1 {
		cnt++
	}
	return cnt
}

func answer() int64 {
	res := (waysRowsAlt() + waysColsAlt() - waysBoth()) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fscan(in, &n, &m, &k)
	rows = make(map[int]*RCData)
	cols = make(map[int]*RCData)
	assignments = make(map[int64]int)
	for i := 0; i < k; i++ {
		var x, y, t int
		fmt.Fscan(in, &x, &y, &t)
		if t == -1 {
			removeCell(x, y)
		} else {
			removeCell(x, y)
			add(x, y, t)
		}
		fmt.Fprintln(out, answer())
	}
}
