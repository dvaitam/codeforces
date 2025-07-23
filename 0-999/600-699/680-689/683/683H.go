package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		a[i]--
	}
	// build direct mapping: friend[u] is the best friend of u
	friend := make([]int, n)
	for i := 0; i < n; i++ {
		friend[a[i]] = i
	}
	visited := make([]bool, n)
	result := make([]int, n)
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		cycle := []int{}
		cur := i
		for !visited[cur] {
			visited[cur] = true
			cycle = append(cycle, cur)
			cur = friend[cur]
		}
		l := len(cycle)
		shift := int(k % int64(l))
		for idx, v := range cycle {
			dest := cycle[(idx+shift)%l]
			result[v] = dest
		}
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for i, v := range result {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v+1)
	}
}
