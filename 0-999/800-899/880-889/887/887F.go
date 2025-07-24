package main

import (
	"bufio"
	"fmt"
	"os"
)

func isNice(a []int, k int) bool {
	n := len(a)
	next := make([]int, n)
	stack := make([]int, 0, n)
	for i := n - 1; i >= 0; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			next[i] = n
		} else {
			next[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}
	for i := 0; i < n; i++ {
		if next[i]-i > k {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	if k >= n {
		fmt.Println("YES")
		return
	}
	if isNice(a, k) {
		fmt.Println("YES")
		return
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if a[i] <= a[j] {
				continue
			}
			a[i], a[j] = a[j], a[i]
			if isNice(a, k) {
				fmt.Println("YES")
				return
			}
			a[i], a[j] = a[j], a[i]
		}
	}
	fmt.Println("NO")
}
