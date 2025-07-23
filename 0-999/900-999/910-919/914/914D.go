package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

var (
	n   int
	arr []int
	st  []int
)

func build(node, l, r int) {
	if l == r {
		st[node] = arr[l]
		return
	}
	mid := (l + r) / 2
	build(node*2, l, mid)
	build(node*2+1, mid+1, r)
	st[node] = gcd(st[node*2], st[node*2+1])
}

func update(node, l, r, pos, val int) {
	if l == r {
		st[node] = val
		arr[pos] = val
		return
	}
	mid := (l + r) / 2
	if pos <= mid {
		update(node*2, l, mid, pos, val)
	} else {
		update(node*2+1, mid+1, r, pos, val)
	}
	st[node] = gcd(st[node*2], st[node*2+1])
}

func query(node, l, r, ql, qr, x int) int {
	if ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		if st[node]%x == 0 {
			return 0
		}
		if l == r {
			return 1
		}
	}
	mid := (l + r) / 2
	res := 0
	if ql <= mid {
		res += query(node*2, l, mid, ql, qr, x)
		if res > 1 {
			return res
		}
	}
	if qr > mid {
		res += query(node*2+1, mid+1, r, ql, qr, x)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n)
	arr = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	st = make([]int, 4*n)
	build(1, 1, n)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var typ int
		fmt.Fscan(reader, &typ)
		if typ == 1 {
			var l, r, x int
			fmt.Fscan(reader, &l, &r, &x)
			cnt := query(1, 1, n, l, r, x)
			if cnt <= 1 {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		} else {
			var idx, val int
			fmt.Fscan(reader, &idx, &val)
			update(1, 1, n, idx, val)
		}
	}
}
