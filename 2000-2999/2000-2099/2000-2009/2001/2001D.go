package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		maxVal := n + 2
		freq := make([]int, maxVal)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			freq[a[i]]++
		}
		visited := make([]bool, maxVal)
		stack := make([]int, 0)
		for _, v := range a {
			freq[v]--
			if visited[v] {
				continue
			}
			for len(stack) > 0 {
				y := stack[len(stack)-1]
				if freq[y] == 0 {
					break
				}
				pos := len(stack)
				better := false
				if pos%2 == 1 {
					if v > y {
						better = true
					}
				} else {
					if v < y {
						better = true
					}
				}
				if !better {
					break
				}
				visited[y] = false
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, v)
			visited[v] = true
		}
		fmt.Fprintln(out, len(stack))
		for i, v := range stack {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
