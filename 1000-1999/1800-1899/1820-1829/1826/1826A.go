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
		res := -1
		for x := 0; x <= n; x++ {
			cnt := 0
			for _, v := range arr {
				if v > x {
					cnt++
				}
			}
			if cnt == x {
				res = x
				break
			}
		}
		fmt.Fprintln(writer, res)
	}
}
