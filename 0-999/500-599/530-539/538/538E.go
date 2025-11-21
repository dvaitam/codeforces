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
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	children := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		children[u] = append(children[u], v)
	}

	leaves := 0
	for i := 0; i < n; i++ {
		if len(children[i]) == 0 {
			leaves++
		}
	}

	order := make([]int, 0, n)
	stack := []int{0}
	parity := make([]int8, n)
	visited := make([]bool, n)
	visited[0] = true

	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		order = append(order, v)
		for _, u := range children[v] {
			if !visited[u] {
				visited[u] = true
				parity[u] = parity[v] ^ 1
				stack = append(stack, u)
			}
		}
	}

	needMax := make([]int, n)
	needMin := make([]int, n)

	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if len(children[v]) == 0 {
			needMax[v] = 1
			needMin[v] = 1
			continue
		}
		if parity[v] == 0 {
			best := int(1e9)
			sum := 0
			for _, u := range children[v] {
				if needMax[u] < best {
					best = needMax[u]
				}
				sum += needMin[u]
			}
			needMax[v] = best
			needMin[v] = sum
		} else {
			sum := 0
			best := int(1e9)
			for _, u := range children[v] {
				sum += needMax[u]
				if needMin[u] < best {
					best = needMin[u]
				}
			}
			needMax[v] = sum
			needMin[v] = best
		}
	}

	maxResult := leaves - needMax[0] + 1
	minResult := needMin[0]

	fmt.Fprintf(out, "%d %d\n", maxResult, minResult)
}
