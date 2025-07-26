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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		b := make([]int, n+1)
		root := -1
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &b[i])
			if b[i] == i {
				root = i
			}
		}
		p := make([]int, n)
		pos := make([]int, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
			pos[p[i]] = i
		}
		if p[0] != root {
			fmt.Fprintln(writer, -1)
			continue
		}
		ok := true
		w := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if i == root {
				w[i] = 0
				continue
			}
			if pos[i] <= pos[b[i]] {
				ok = false
				break
			}
			w[i] = pos[i] - pos[b[i]]
		}
		if !ok {
			fmt.Fprintln(writer, -1)
			continue
		}
		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, w[i])
		}
		fmt.Fprintln(writer)
	}
}
