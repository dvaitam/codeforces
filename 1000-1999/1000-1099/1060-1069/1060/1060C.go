package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	var X int64
	fmt.Fscan(reader, &n, &m)
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		var v int64
		fmt.Fscan(reader, &v)
		a[i] = a[i-1] + v
	}
	b := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		var v int64
		fmt.Fscan(reader, &v)
		b[i] = b[i-1] + v
	}
	fmt.Fscan(reader, &X)

	f := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		f[i] = 1 << 60
		for j := i; j <= n; j++ {
			sum := a[j] - a[j-i]
			if sum < f[i] {
				f[i] = sum
			}
		}
	}

	g := make([]int64, m+1)
	for i := 1; i <= m; i++ {
		g[i] = 1 << 60
		for j := i; j <= m; j++ {
			sum := b[j] - b[j-i]
			if sum < g[i] {
				g[i] = sum
			}
		}
	}

	ans := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if f[i]*g[j] <= X {
				prod := i * j
				if prod > ans {
					ans = prod
				}
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
