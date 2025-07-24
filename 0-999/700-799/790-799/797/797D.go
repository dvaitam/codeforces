package main

import (
	"bufio"
	"fmt"
	"os"
)

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	val := make([]int, n+1)
	left := make([]int, n+1)
	right := make([]int, n+1)
	child := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &val[i], &left[i], &right[i])
		if left[i] != -1 {
			child[left[i]] = true
		}
		if right[i] != -1 {
			child[right[i]] = true
		}
	}
	root := 1
	for i := 1; i <= n; i++ {
		if !child[i] {
			root = i
			break
		}
	}

	type item struct {
		idx       int
		low, high int64
	}
	stack := []item{{root, -1 << 62, 1 << 62}}
	found := make(map[int]bool)
	freq := make(map[int]int)

	for len(stack) > 0 {
		it := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		v := val[it.idx]
		freq[v]++
		if int64(v) > it.low && int64(v) < it.high {
			found[v] = true
		}
		if left[it.idx] != -1 {
			nhigh := minInt64(it.high, int64(v))
			stack = append(stack, item{left[it.idx], it.low, nhigh})
		}
		if right[it.idx] != -1 {
			nlow := maxInt64(it.low, int64(v))
			stack = append(stack, item{right[it.idx], nlow, it.high})
		}
	}

	success := 0
	for val, ok := range found {
		if ok {
			success += freq[val]
		}
	}
	fmt.Println(n - success)
}
