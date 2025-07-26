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
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
			b[i]--
		}
		if k == 1 {
			ok := true
			for i := 0; i < n; i++ {
				if b[i] != i {
					ok = false
					break
				}
			}
			if ok {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
			continue
		}
		visited := make([]int, n)
		ok := true
		for i := 0; i < n && ok; i++ {
			if visited[i] != 0 {
				continue
			}
			x := i
			for visited[x] == 0 {
				visited[x] = 1
				x = b[x]
			}
			if visited[x] == 1 {
				// found a new cycle starting at x
				cnt := 1
				y := b[x]
				for y != x {
					cnt++
					y = b[y]
				}
				if cnt != k {
					ok = false
					break
				}
			}
			// mark the nodes on the path as processed
			x = i
			for visited[x] == 1 {
				visited[x] = 2
				x = b[x]
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
