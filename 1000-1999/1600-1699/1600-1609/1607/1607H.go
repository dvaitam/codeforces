package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		type Dish struct {
			idx     int
			a, b, m int
			c       int
			l, r    int
		}
		groups := make(map[int][]Dish)
		ansA := make([]int, n)
		for i := 0; i < n; i++ {
			var a, b, m int
			fmt.Fscan(in, &a, &b, &m)
			c := a + b - m
			l := c - b
			if l < 0 {
				l = 0
			}
			r := c
			if r > a {
				r = a
			}
			groups[c] = append(groups[c], Dish{idx: i, a: a, b: b, m: m, c: c, l: l, r: r})
		}
		variety := 0
		for c, arr := range groups {
			sort.Slice(arr, func(i, j int) bool {
				if arr[i].r == arr[j].r {
					return arr[i].l < arr[j].l
				}
				return arr[i].r < arr[j].r
			})
			last := -1
			for _, d := range arr {
				if last < d.l {
					last = d.r
					variety++
				}
				ansA[d.idx] = last
			}
			groups[c] = arr
		}
		fmt.Fprintln(out, variety)
		for c, arr := range groups {
			_ = c
			for _, d := range arr {
				x := d.a - ansA[d.idx]
				y := d.m - x
				fmt.Fprintln(out, x, y)
			}
		}
	}
}
