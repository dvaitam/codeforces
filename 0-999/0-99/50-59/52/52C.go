package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const INF = int64(1e18)

var t []int64
var lazy []int64
var a []int64

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func build(node, l, r int) {
	if l == r {
		t[node] = a[l]
		return
	}
	m := (l + r) >> 1
	build(2*node, l, m)
	build(2*node+1, m+1, r)
	t[node] = min(t[2*node], t[2*node+1])
}

func push(node int) {
	v := lazy[node]
	if v != 0 {
		// propagate to children
		t[2*node] += v
		lazy[2*node] += v
		t[2*node+1] += v
		lazy[2*node+1] += v
		lazy[node] = 0
	}
}

func update(node, l, r, ql, qr int, v int64) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		t[node] += v
		lazy[node] += v
		return
	}
	push(node)
	m := (l + r) >> 1
	update(2*node, l, m, ql, qr, v)
	update(2*node+1, m+1, r, ql, qr, v)
	t[node] = min(t[2*node], t[2*node+1])
}

func query(node, l, r, ql, qr int) int64 {
	if ql > r || qr < l {
		return INF
	}
	if ql <= l && r <= qr {
		return t[node]
	}
	push(node)
	m := (l + r) >> 1
	left := query(2*node, l, m, ql, qr)
	right := query(2*node+1, m+1, r, ql, qr)
	return min(left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a = make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	t = make([]int64, 4*n)
	lazy = make([]int64, 4*n)
	build(1, 0, n-1)

	var m int
	fmt.Fscan(reader, &m)
	// consume end of line
	reader.ReadString('\n')
	for i := 0; i < m; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && len(line) == 0 {
			break
		}
		fields := strings.Fields(line)
		lf, _ := strconv.Atoi(fields[0])
		rg, _ := strconv.Atoi(fields[1])
		if len(fields) == 2 {
			var res int64
			if lf <= rg {
				res = query(1, 0, n-1, lf, rg)
			} else {
				res = min(
					query(1, 0, n-1, lf, n-1),
					query(1, 0, n-1, 0, rg),
				)
			}
			fmt.Fprintln(writer, res)
		} else {
			v, _ := strconv.ParseInt(fields[2], 10, 64)
			if lf <= rg {
				update(1, 0, n-1, lf, rg, v)
			} else {
				update(1, 0, n-1, lf, n-1, v)
				update(1, 0, n-1, 0, rg, v)
			}
		}
	}
}
