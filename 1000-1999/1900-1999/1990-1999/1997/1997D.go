package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 200005

var a [N]int
var e [N][]int

func S(u int) int {
	if len(e[u]) == 0 {
		return a[u]
	}
	r := 1_000_000_000
	for _, v := range e[u] {
		val := S(v)
		if val < r {
			r = val
		}
	}
	if a[u] < r && u != 1 {
		return (r + a[u]) / 2
	}
	return r
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t, n int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		fmt.Fscan(reader, &n)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			e[i] = e[i][:0]
		}
		for i := 2; i <= n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			e[x] = append(e[x], i)
		}
		res := a[1] + S(1)
		fmt.Fprintln(writer, res)
	}
}
