package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	h := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &h[i])
	}
	if n == 1 {
		fmt.Println(h[0])
		return
	}
	left := make([]int, n)
	right := make([]int, n)
	left[0] = 1
	if left[0] > h[0] {
		left[0] = h[0]
	}
	for i := 1; i < n; i++ {
		left[i] = left[i-1] + 1
		if left[i] > h[i] {
			left[i] = h[i]
		}
	}
	right[n-1] = 1
	if right[n-1] > h[n-1] {
		right[n-1] = h[n-1]
	}
	for i := n - 2; i >= 0; i-- {
		right[i] = right[i+1] + 1
		if right[i] > h[i] {
			right[i] = h[i]
		}
	}
	ans := 0
	for i := 0; i < n; i++ {
		v := left[i]
		if right[i] < v {
			v = right[i]
		}
		if v > ans {
			ans = v
		}
	}
	fmt.Println(ans)
}
