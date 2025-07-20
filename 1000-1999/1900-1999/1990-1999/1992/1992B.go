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
		maxVal := 0
		var ai int
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &ai)
			if ai > maxVal {
				maxVal = ai
			}
		}
		ans := 2*(n-maxVal) - (k - 1)
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(writer, ans)
	}
}
