package main

import (
	"bufio"
	"fmt"
	"os"
)

type void struct{}

var (
	qd  [105][105]bool
	ans [105][105]int
	n   int
	in  = bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func query(x, y int) int {
	if qd[x][y] {
		return ans[x][y]
	}
	qd[x][y] = true
	fmt.Fprintf(out, "? %d %d\n", x+1, y)
	out.Flush()
	var g int
	fmt.Fscan(in, &g)
	ans[x][y] = g
	return g
}

func sol(v []int, x int) {
	if x == -1 || len(v) <= 1 {
		return
	}
	a := v[0]
	b := v[len(v)-1]
	df := query(a, b) == x
	fa := []int{a}
	fb := []int{}
	for i := 1; i < len(v)-1; i++ {
		j := v[i]
		var hi bool
		if abs(j-a) <= abs(j-b) {
			hi = query(j, b) == x
			if df {
				hi = !hi
			}
		} else {
			hi = query(a, j) == x
		}
		if hi {
			fb = append(fb, j)
		} else {
			fa = append(fa, j)
		}
	}
	if df {
		fb = append(fb, b)
	} else {
		fa = append(fa, b)
	}
	for _, i := range fa {
		for _, j := range fb {
			ans[i][j] = x
			ans[j][i] = x
		}
	}
	sol(fa, x-1)
	sol(fb, x-1)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solve() {
	fmt.Fscan(in, &n)
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			qd[i][j] = false
			ans[i][j] = -1
		}
	}
	v := make([]int, n+1)
	for i := 0; i <= n; i++ {
		v[i] = i
	}
	sol(v, 29)
	fmt.Fprintln(out, "!")
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			fmt.Fprintf(out, "%d ", ans[i][j])
		}
		fmt.Fprintln(out)
	}
	out.Flush()
}

func main() {
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		solve()
	}
}

