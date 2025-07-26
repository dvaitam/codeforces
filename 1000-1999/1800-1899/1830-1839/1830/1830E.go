package main

import (
	"bufio"
	"fmt"
	"os"
)

func bullyCount(p []int) int {
	n := len(p)
	arr := make([]int, n)
	copy(arr, p)
	steps := 0
	for {
		i := -1
		maxVal := -1
		for idx, val := range arr {
			if val != idx+1 && val > maxVal {
				maxVal = val
				i = idx
			}
		}
		if i == -1 {
			break
		}
		j := -1
		minVal := int(1e9)
		for k := i + 1; k < n; k++ {
			if arr[k] < minVal {
				minVal = arr[k]
				j = k
			}
		}
		arr[i], arr[j] = arr[j], arr[i]
		steps++
	}
	return steps
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}
	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
	}
	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		p[x], p[y] = p[y], p[x]
		fmt.Fprintln(out, bullyCount(p))
	}
}
