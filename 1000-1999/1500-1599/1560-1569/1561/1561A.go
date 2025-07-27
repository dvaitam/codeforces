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
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		ans := 0
		for {
			if isSorted(arr) {
				break
			}
			if ans%2 == 0 { // iteration number is odd
				for j := 0; j+1 < n; j += 2 {
					if arr[j] > arr[j+1] {
						arr[j], arr[j+1] = arr[j+1], arr[j]
					}
				}
			} else { // iteration number is even
				for j := 1; j+1 < n; j += 2 {
					if arr[j] > arr[j+1] {
						arr[j], arr[j+1] = arr[j+1], arr[j]
					}
				}
			}
			ans++
		}
		fmt.Fprintln(out, ans)
	}
}
