package main

import (
	"bufio"
	"fmt"
	"os"
)

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
		minVal := 101
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			if arr[i] < minVal {
				minVal = arr[i]
			}
		}
		cnt := 0
		for _, v := range arr {
			if v > minVal {
				cnt++
			}
		}
		fmt.Fprintln(out, cnt)
	}
}
