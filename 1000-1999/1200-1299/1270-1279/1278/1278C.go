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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Fscan(reader, &arr[i])
			if arr[i] == 1 {
				arr[i] = 1
			} else {
				arr[i] = -1
			}
		}
		total := 0
		for _, v := range arr {
			total += v
		}
		right := map[int]int{0: 0}
		diff := 0
		for i := 0; i < n; i++ {
			diff += arr[n+i]
			if _, ok := right[diff]; !ok {
				right[diff] = i + 1
			}
		}
		ans := 2 * n
		if val, ok := right[total]; ok && val < ans {
			ans = val
		}
		diffLeft := 0
		for j := 1; j <= n; j++ {
			diffLeft += arr[n-j]
			need := total - diffLeft
			if k, ok := right[need]; ok {
				if j+k < ans {
					ans = j + k
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
