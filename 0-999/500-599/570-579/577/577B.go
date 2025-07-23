package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &nums[i])
	}
	if n >= m {
		fmt.Println("YES")
		return
	}
	reachable := make([]bool, m)
	for _, v := range nums {
		v %= m
		next := make([]bool, m)
		copy(next, reachable)
		next[v] = true
		for r := 0; r < m; r++ {
			if reachable[r] {
				next[(r+v)%m] = true
			}
		}
		reachable = next
		if reachable[0] {
			fmt.Println("YES")
			return
		}
	}
	if reachable[0] {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
