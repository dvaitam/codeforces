package main

import (
	"bufio"
	"fmt"
	"os"
)

func isSorted(a []int) bool {
	for i := 0; i+1 < len(a); i++ {
		if a[i] > a[i+1] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		ans := 0
		for !isSorted(arr) {
			if ans%2 == 0 { // odd iteration number (1-based)
				for i := 0; i+1 < n; i += 2 {
					if arr[i] > arr[i+1] {
						arr[i], arr[i+1] = arr[i+1], arr[i]
					}
				}
			} else {
				for i := 1; i+1 < n; i += 2 {
					if arr[i] > arr[i+1] {
						arr[i], arr[i+1] = arr[i+1], arr[i]
					}
				}
			}
			ans++
		}
		fmt.Fprintln(writer, ans)
	}
}
