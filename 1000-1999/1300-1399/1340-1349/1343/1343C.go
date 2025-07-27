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
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		if n == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		currMax := arr[0]
		ans := int64(0)
		currSign := arr[0] > 0
		for i := 1; i < n; i++ {
			sign := arr[i] > 0
			if sign == currSign {
				if arr[i] > currMax {
					currMax = arr[i]
				}
			} else {
				ans += currMax
				currMax = arr[i]
				currSign = sign
			}
		}
		ans += currMax
		fmt.Fprintln(writer, ans)
	}
}
