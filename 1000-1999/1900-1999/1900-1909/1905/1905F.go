package main

import (
	"bufio"
	"fmt"
	"os"
)

func good(p []int) int {
	n := len(p)
	prefixMax := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		if p[i] > maxVal {
			maxVal = p[i]
		}
		prefixMax[i] = maxVal
	}
	suffixMin := make([]int, n)
	minVal := n + 1
	for i := n - 1; i >= 0; i-- {
		if p[i] < minVal {
			minVal = p[i]
		}
		suffixMin[i] = minVal
	}
	cnt := 0
	for i := 0; i < n; i++ {
		leftMax := 0
		if i > 0 {
			leftMax = prefixMax[i-1]
		}
		rightMin := n + 1
		if i+1 < n {
			rightMin = suffixMin[i+1]
		}
		if leftMax < p[i] && p[i] < rightMin {
			cnt++
		}
	}
	return cnt
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(reader, &t)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		best := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				p[i], p[j] = p[j], p[i]
				g := good(p)
				if g > best {
					best = g
				}
				p[i], p[j] = p[j], p[i]
			}
		}
		fmt.Fprintln(writer, best)
	}
}
