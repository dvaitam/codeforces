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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		minVal, maxVal := a[0], a[0]
		for _, v := range a {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
			}
		}
		if minVal == maxVal {
			ans := int64(n) * int64(n-1)
			fmt.Fprintln(writer, ans)
			continue
		}
		var cntMin, cntMax int64
		for _, v := range a {
			if v == minVal {
				cntMin++
			}
			if v == maxVal {
				cntMax++
			}
		}
		ans := cntMin * cntMax * 2
		fmt.Fprintln(writer, ans)
	}
}
