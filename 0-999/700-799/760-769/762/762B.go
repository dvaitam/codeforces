package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b, c int
	if _, err := fmt.Fscan(in, &a, &b, &c); err != nil {
		return
	}
	var m int
	fmt.Fscan(in, &m)
	usb := make([]int64, 0)
	ps := make([]int64, 0)
	for i := 0; i < m; i++ {
		var cost int64
		var typ string
		fmt.Fscan(in, &cost, &typ)
		if typ == "USB" {
			usb = append(usb, cost)
		} else {
			ps = append(ps, cost)
		}
	}
	sort.Slice(usb, func(i, j int) bool { return usb[i] < usb[j] })
	sort.Slice(ps, func(i, j int) bool { return ps[i] < ps[j] })
	var count int
	var total int64
	// use USB mouses for USB-only computers
	u := 0
	for u < len(usb) && a > 0 {
		total += usb[u]
		count++
		u++
		a--
	}
	// use PS/2 mouses for PS/2-only computers
	p := 0
	for p < len(ps) && b > 0 {
		total += ps[p]
		count++
		p++
		b--
	}
	// collect remaining mouses
	rest := make([]int64, 0, len(usb)-u+len(ps)-p)
	rest = append(rest, usb[u:]...)
	rest = append(rest, ps[p:]...)
	sort.Slice(rest, func(i, j int) bool { return rest[i] < rest[j] })
	r := 0
	for r < len(rest) && c > 0 {
		total += rest[r]
		count++
		r++
		c--
	}
	fmt.Printf("%d %d\n", count, total)
}
