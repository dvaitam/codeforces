package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	sizes := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &sizes[i])
	}

	maxVal := 0
	if n <= k {
		maxVal = sizes[n-1]
	} else {
		single := 2*k - n
		for i := 0; i < single; i++ {
			if sizes[i] > maxVal {
				maxVal = sizes[i]
			}
		}
		left, right := single, n-1
		for left < right {
			sum := sizes[left] + sizes[right]
			if sum > maxVal {
				maxVal = sum
			}
			left++
			right--
		}
	}

	fmt.Fprintln(writer, maxVal)
}
