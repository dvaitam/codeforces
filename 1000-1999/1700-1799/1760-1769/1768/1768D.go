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
		var n int
		fmt.Fscan(reader, &n)
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		visited := make([]bool, n+1)
		cycleID := make([]int, n+1)
		cid := 0
		for i := 1; i <= n; i++ {
			if !visited[i] {
				cid++
				j := i
				for !visited[j] {
					visited[j] = true
					cycleID[j] = cid
					j = p[j]
				}
			}
		}
		base := n - cid
		same := false
		for i := 1; i < n; i++ {
			if cycleID[i] == cycleID[i+1] {
				same = true
				break
			}
		}
		ans := base + 1
		if same {
			ans = base - 1
		}
		fmt.Fprintln(writer, ans)
	}
}
