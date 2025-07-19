package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n      int
	v      [][]int
	o      []int
	tam    []int
	ans    []int
	totSum int64
)

func findSize(a, p int) {
	tam[a] = 1
	o = append(o, a)
	for _, x := range v[a] {
		if x == p {
			continue
		}
		findSize(x, a)
		// accumulate contribution
		var smaller = tam[x]
		if n-tam[x] < smaller {
			smaller = n - tam[x]
		}
		totSum += int64(smaller)
		tam[a] += tam[x]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	fmt.Fscan(in, &n)
	v = make([][]int, n+1)
	tam = make([]int, n+1)
	ans = make([]int, n+1)
	o = make([]int, 0, n)

	for i := 1; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		v[a] = append(v[a], b)
		v[b] = append(v[b], a)
	}

	findSize(1, 0)

	half := n / 2
	for i := 0; i < n; i++ {
		a := o[i]
		b := o[(i+half)%n]
		ans[b] = a
	}

	// total sum doubled
	fmt.Fprintln(out, totSum*2)
	for i := 1; i <= n; i++ {
		if i > 1 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
