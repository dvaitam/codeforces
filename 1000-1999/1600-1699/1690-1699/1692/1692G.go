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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}

		cnt := 0
		ans := 0
		for i := 1; i < n; i++ {
			if arr[i-1] < arr[i]*2 {
				cnt++
			} else {
				cnt = 0
			}
			if cnt >= k {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
