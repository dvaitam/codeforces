package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Edge struct {
	to   int
	aVal int64
	bVal int64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		g := make([][]Edge, n+1)
		for i := 2; i <= n; i++ {
			var p int
			var a, b int64
			fmt.Fscan(reader, &p, &a, &b)
			g[p] = append(g[p], Edge{to: i, aVal: a, bVal: b})
		}

		res := make([]int, n+1)
		prefix := []int64{0}
		type Node struct {
			id   int
			aSum int64
			idx  int
		}
		stack := []Node{{id: 1, aSum: 0, idx: 0}}

		for len(stack) > 0 {
			top := &stack[len(stack)-1]
			if top.idx == len(g[top.id]) {
				if top.id != 1 {
					prefix = prefix[:len(prefix)-1]
				}
				stack = stack[:len(stack)-1]
				continue
			}
			e := g[top.id][top.idx]
			top.idx++
			newASum := top.aSum + e.aVal
			prefix = append(prefix, prefix[len(prefix)-1]+e.bVal)
			pos := sort.Search(len(prefix), func(i int) bool {
				return prefix[i] > newASum
			}) - 1
			res[e.to] = pos
			stack = append(stack, Node{id: e.to, aSum: newASum, idx: 0})
		}

		for i := 2; i <= n; i++ {
			if i > 2 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, res[i])
		}
		fmt.Fprintln(writer)
	}
}
