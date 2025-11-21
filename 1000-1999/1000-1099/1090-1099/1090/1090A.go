package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	maxPerCompany := make([]int64, n)
	size := make([]int64, n)
	var globalMax int64

	for i := 0; i < n; i++ {
		var m int
		if _, err := fmt.Fscan(in, &m); err != nil {
			return
		}
		size[i] = int64(m)
		var localMax int64
		for j := 0; j < m; j++ {
			var salary int64
			fmt.Fscan(in, &salary)
			if salary > localMax {
				localMax = salary
			}
		}
		maxPerCompany[i] = localMax
		if localMax > globalMax {
			globalMax = localMax
		}
	}

	var totalIncrease int64
	for i := 0; i < n; i++ {
		totalIncrease += (globalMax - maxPerCompany[i]) * size[i]
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, totalIncrease)
}
