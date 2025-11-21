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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for test := 0; test < t; test++ {
		var n int
		fmt.Fscan(reader, &n)

		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &p[i])
		}

		var s string
		fmt.Fscan(reader, &s)

		colorBlack := make([]bool, n+1)
		for i := 1; i <= n; i++ {
			val := p[i]
			colorBlack[val] = s[i-1] == '0'
		}

		visited := make([]bool, n+1)
		ans := make([]int, n+1)

		for i := 1; i <= n; i++ {
			if visited[i] {
				continue
			}
			cur := i
			cycle := make([]int, 0)
			for !visited[cur] {
				visited[cur] = true
				cycle = append(cycle, cur)
				cur = p[cur]
			}

			countBlack := 0
			for _, node := range cycle {
				if colorBlack[node] {
					countBlack++
				}
			}
			for _, node := range cycle {
				ans[node] = countBlack
			}
		}

		for i := 1; i <= n; i++ {
			if i > 1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans[i])
		}
		if test != t-1 {
			fmt.Fprintln(writer)
		}
	}
}
