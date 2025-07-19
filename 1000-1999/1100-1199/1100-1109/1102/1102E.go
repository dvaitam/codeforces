package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod = 998244353

type Node struct {
	id, val int
}

type Node1 struct {
	l, r, val int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n int
	fmt.Fscan(reader, &n)
	a := make([]Node, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i].val)
		a[i].id = i
	}
	sort.Slice(a, func(i, j int) bool {
		if a[i].val != a[j].val {
			return a[i].val < a[j].val
		}
		return a[i].id < a[j].id
	})
	b := make([]Node1, 0)
	for i := 1; i < n; i++ {
		if a[i].val == a[i-1].val {
			b = append(b, Node1{a[i-1].id, a[i].id, a[i].val})
		}
	}
	cnt := len(b)
	sort.Slice(b, func(i, j int) bool {
		return b[i].l < b[j].l
	})
	ans := make([]int, n)
	if cnt > 0 {
		lp, ip, col := b[0].l, 0, 1
		for lp < n && ip < cnt {
			for lp > b[ip].r {
				ip++
				if ip >= cnt {
					break
				}
				if b[ip].val != b[ip-1].val && lp <= b[ip].l {
					col++
				}
				if lp < b[ip].l {
					lp = b[ip].l
				}
			}
			if ip >= cnt {
				break
			}
			ans[lp] = col
			lp++
		}
	}
	var anss int64 = 1
	for i := 1; i < n; i++ {
		if ans[i] == 0 || ans[i] != ans[i-1] {
			anss = anss * 2 % mod
		}
	}
	fmt.Fprintln(writer, anss)
}
