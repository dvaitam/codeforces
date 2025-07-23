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
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	const maxVal = 1000000 + 1
	count := make([]int, maxVal)

	left := 0
	distinct := 0
	bestL, bestR := 0, 0

	for right, v := range arr {
		if count[v] == 0 {
			distinct++
		}
		count[v]++
		for distinct > k {
			c := arr[left]
			count[c]--
			if count[c] == 0 {
				distinct--
			}
			left++
		}
		if right-left > bestR-bestL {
			bestL = left
			bestR = right
		}
	}

	fmt.Fprintln(writer, bestL+1, bestR+1)
}
