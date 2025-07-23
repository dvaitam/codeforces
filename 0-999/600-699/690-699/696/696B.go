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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	parents := make([]int, n+1)
	children := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		var p int
		fmt.Fscan(reader, &p)
		parents[i] = p
		children[p] = append(children[p], i)
	}

	order := make([]int, 0, n)
	stack := []int{1}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, c := range children[v] {
			stack = append(stack, c)
		}
	}

	size := make([]int, n+1)
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		s := 1
		for _, c := range children[v] {
			s += size[c]
		}
		size[v] = s
	}

	ans := make([]float64, n+1)
	ans[1] = 1.0
	for _, v := range order {
		for _, c := range children[v] {
			ans[c] = ans[v] + 1.0 + float64(size[v]-1-size[c])/2.0
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprintf(writer, "%.10f", ans[i])
	}
	fmt.Fprintln(writer)
}
