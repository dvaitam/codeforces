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
		var s string
		fmt.Fscan(reader, &s)
		l := make([]int, n+1)
		r := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &l[i], &r[i])
		}
		type node struct{ idx, cost int }
		stack := []node{{1, 0}}
		ans := int(1e9)
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			i := cur.idx
			if l[i] == 0 && r[i] == 0 {
				if cur.cost < ans {
					ans = cur.cost
				}
			}
			if l[i] != 0 {
				c := cur.cost
				if s[i-1] != 'L' {
					c++
				}
				stack = append(stack, node{l[i], c})
			}
			if r[i] != 0 {
				c := cur.cost
				if s[i-1] != 'R' {
					c++
				}
				stack = append(stack, node{r[i], c})
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
