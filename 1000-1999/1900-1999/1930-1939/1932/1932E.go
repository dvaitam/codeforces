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
		var s string
		fmt.Fscan(reader, &s)

		diff := make([]int64, n+1)
		for i := 0; i < n; i++ {
			d := int64(s[n-1-i] - '0')
			diff[0] += d
			diff[i+1] -= d
		}

		arr := make([]int64, n+1)
		arr[0] = diff[0]
		for i := 1; i <= n; i++ {
			arr[i] = arr[i-1] + diff[i]
		}

		for i := 0; i < len(arr); i++ {
			carry := arr[i] / 10
			arr[i] %= 10
			if carry > 0 {
				if i+1 == len(arr) {
					arr = append(arr, 0)
				}
				arr[i+1] += carry
			}
		}

		idx := len(arr) - 1
		for idx > 0 && arr[idx] == 0 {
			idx--
		}
		for i := idx; i >= 0; i-- {
			fmt.Fprint(writer, arr[i])
		}
		fmt.Fprintln(writer)
	}
}
