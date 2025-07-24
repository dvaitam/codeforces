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
		arr := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			sum += arr[i]
		}
		if sum%2 == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		minOps := int(1e9)
		for _, x := range arr {
			cnt := 0
			y := x
			if y%2 == 0 {
				for y%2 == 0 {
					y /= 2
					cnt++
				}
			} else {
				for y%2 == 1 {
					y /= 2
					cnt++
				}
			}
			if cnt < minOps {
				minOps = cnt
			}
		}
		fmt.Fprintln(writer, minOps)
	}
}
