package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func calc(x, y, z int64) int64 {
	a := x - y
	b := y - z
	c := z - x
	return a*a + b*b + c*c
}

func best(a, b, c []int64) int64 {
	res := int64(1<<63 - 1)
	for _, y := range b {
		// x <= y
		xi := sort.Search(len(a), func(i int) bool { return a[i] > y }) - 1
		if xi < 0 {
			continue
		}
		// z >= y
		zi := sort.Search(len(c), func(i int) bool { return c[i] >= y })
		if zi == len(c) {
			continue
		}
		v := calc(a[xi], y, c[zi])
		if v < res {
			res = v
		}
	}
	return res
}

func solve(r, g, b []int64) int64 {
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
	sort.Slice(g, func(i, j int) bool { return g[i] < g[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })

	res := int64(1<<63 - 1)
	res = min(res, best(r, g, b))
	res = min(res, best(r, b, g))
	res = min(res, best(g, r, b))
	res = min(res, best(g, b, r))
	res = min(res, best(b, r, g))
	res = min(res, best(b, g, r))
	return res
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var nr, ng, nb int
		fmt.Fscan(reader, &nr, &ng, &nb)
		r := make([]int64, nr)
		gArr := make([]int64, ng)
		bArr := make([]int64, nb)
		for i := 0; i < nr; i++ {
			fmt.Fscan(reader, &r[i])
		}
		for i := 0; i < ng; i++ {
			fmt.Fscan(reader, &gArr[i])
		}
		for i := 0; i < nb; i++ {
			fmt.Fscan(reader, &bArr[i])
		}
		ans := solve(r, gArr, bArr)
		fmt.Fprintln(writer, ans)
	}
}
