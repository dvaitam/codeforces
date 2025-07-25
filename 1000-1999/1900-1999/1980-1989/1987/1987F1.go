package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxOperations(a []int) int {
	memo := make(map[string]int)
	var dfs func([]int) int
	dfs = func(arr []int) int {
		if len(arr) < 2 {
			return 0
		}
		key := fmt.Sprint(arr)
		if v, ok := memo[key]; ok {
			return v
		}
		best := 0
		for i := 0; i < len(arr)-1; i++ {
			if arr[i] == i+1 {
				b := append([]int{}, arr[:i]...)
				b = append(b, arr[i+2:]...)
				if val := 1 + dfs(b); val > best {
					best = val
				}
			}
		}
		memo[key] = best
		return best
	}
	return dfs(a)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		fmt.Fprintln(writer, maxOperations(arr))
	}
}
