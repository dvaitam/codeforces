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
	var t int
	fmt.Fscan(reader, &t)
	for tt := 0; tt < t; tt++ {
		solve(reader, writer)
	}
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
	var n, m int
	fmt.Fscan(reader, &n, &m)
	u := make([]int, m)
	v := make([]int, m)
	deg := make([]int, n+1)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &u[i], &v[i])
		deg[u[i]]++
		deg[v[i]]++
	}
	var foundX int
	for x := 1; x <= n; x++ {
		if deg[x] < n-1 {
			foundX = x
			break
		}
	}
	if foundX != 0 {
		fmt.Fprintln(writer, 2)
		for i := 0; i < m; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			if u[i] == foundX || v[i] == foundX {
				fmt.Fprint(writer, 1)
			} else {
				fmt.Fprint(writer, 2)
			}
		}
		writer.WriteByte('\n')
		return
	}
	fmt.Fprintln(writer, 3)
	for i := 0; i < m; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		var x int
		if u[i]+v[i] == 3 {
			x = 1
		} else if u[i] == 1 || v[i] == 1 {
			x = 2
		} else {
			x = 3
		}
		fmt.Fprint(writer, x)
	}
	writer.WriteByte('\n')
}
