package main

import (
	"bufio"
	"fmt"
	"os"
)

func reachable(start, target int, edges [][]bool) bool {
	n := 20
	visited := make([]bool, n)
	stack := []int{start}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if v == target {
			return true
		}
		if visited[v] {
			continue
		}
		visited[v] = true
		for i := 0; i < n; i++ {
			if edges[v][i] && !visited[i] {
				stack = append(stack, i)
			}
		}
	}
	return false
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
		var A, B string
		fmt.Fscan(reader, &A)
		fmt.Fscan(reader, &B)

		edges := make([][]bool, 20)
		for i := range edges {
			edges[i] = make([]bool, 20)
		}
		for i := 0; i < n; i++ {
			x := int(A[i] - 'a')
			y := int(B[i] - 'a')
			if x != y {
				edges[x][y] = true
			}
		}

		count := 0
		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				if edges[i][j] {
					count++
				}
			}
		}

		for i := 0; i < 20; i++ {
			for j := 0; j < 20; j++ {
				if edges[i][j] {
					edges[i][j] = false
					if reachable(i, j, edges) {
						count--
					} else {
						edges[i][j] = true
					}
				}
			}
		}

		fmt.Fprintln(writer, count)
	}
}
