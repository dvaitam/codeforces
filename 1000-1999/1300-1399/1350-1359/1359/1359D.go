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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	ans := 0
	for m := -30; m <= 30; m++ {
		sum := 0
		hasMax := false
		for i := 0; i < n; i++ {
			if arr[i] > m {
				sum = 0
				hasMax = false
				continue
			}
			sum += arr[i]
			if arr[i] == m {
				hasMax = true
			}
			if hasMax {
				val := sum - m
				if val > ans {
					ans = val
				}
			}
			if sum < 0 {
				sum = 0
				hasMax = false
			}
		}
	}

	fmt.Fprintln(writer, ans)
}
