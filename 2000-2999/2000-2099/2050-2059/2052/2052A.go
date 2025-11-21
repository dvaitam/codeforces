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

	var n int
	fmt.Fscan(in, &n)
	final := make([]int, n)
	pos := make([]int, n+1)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &final[i])
		pos[final[i]] = i
	}

	type pair struct{ x, y int }
	var ops []pair

	current := make([]int, n)
	for i := 0; i < n; i++ {
		current[i] = i + 1
	}
	index := make([]int, n+1)
	for i := 0; i < n; i++ {
		index[current[i]] = i
	}

	for i := n - 1; i >= 0; i-- {
		car := final[i]
		for index[car] < i {
			j := index[car]
			ops = append(ops, pair{current[j+1], current[j]})
			current[j], current[j+1] = current[j+1], current[j]
			index[current[j]] = j
			index[current[j+1]] = j + 1
		}
	}

	fmt.Fprintln(out, len(ops))
	for _, op := range ops {
		fmt.Fprintln(out, op.x, op.y)
	}
}
