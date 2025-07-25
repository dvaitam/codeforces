package main

import (
	"bufio"
	"fmt"
	"os"
)

func compute(a, b []int64, c []int64) int64 {
	n := len(a)
	var wine int64
	carry := int64(0)
	for i := 0; i < n; i++ {
		water := a[i] + carry
		removed := b[i]
		if removed > water {
			removed = water
		}
		wine += removed
		water -= removed
		carry = water
		if i < n-1 {
			if carry > c[i] {
				carry = c[i]
			}
		}
	}
	return wine
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	a := make([]int64, n)
	b := make([]int64, n)
	c := make([]int64, n-1)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	for i := 0; i < n-1; i++ {
		fmt.Fscan(reader, &c[i])
	}

	for ; q > 0; q-- {
		var p int
		var x, y, z int64
		fmt.Fscan(reader, &p, &x, &y, &z)
		p--
		a[p] = x
		b[p] = y
		if p < n-1 {
			c[p] = z
		}
		result := compute(a, b, c)
		fmt.Fprintln(writer, result)
	}
}
