package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const _ = 500005

func p(l, r int64, a, b, c []int64) bool {
	n := r - l + 1
	if n < 6 {
		return false
	}
	for i := int64(0); i < n; i++ {
		b[i] = a[i+l]
	}
	sort.Slice(b[:n], func(i, j int) bool { return b[i] < b[j] })
	for i := int64(2); i < n; i++ {
		t, j := int64(0), int64(0)
		for ; j < i-1; j++ {
			if b[j]+b[i-1] > b[i] {
				break
			}
		}
		if j == i-1 {
			continue
		}
		for j++; j < n; j++ {
			if j != i && j != i-1 {
				c[t] = b[j]
				t++
			}
		}
		for j := int64(2); j < t; j++ {
			if c[j-2]+c[j-1] > c[j] {
				return true
			}
		}
	}
	return false
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q, x, y int64
	fmt.Fscan(in, &n, &q)
	a := make([]int64, n+1)
	b := make([]int64, n+1)
	c := make([]int64, n+1)
	d := make([]int64, n+1)
	for i := int64(1); i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}
	for i, j := int64(1), int64(1); i <= n; i++ {
		for ; j <= n && !p(i, j, a, b, c); j++ {
		}
		d[i] = j
	}
	for ; q > 0; q-- {
		fmt.Fscan(in, &x, &y)
		if y < d[x] {
			fmt.Println("NO")
		} else {
			fmt.Println("YES")
		}
	}
}
