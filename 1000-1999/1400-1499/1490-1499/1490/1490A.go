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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		ans := 0
		for i := 0; i < n-1; i++ {
			a := arr[i]
			b := arr[i+1]
			if a < b {
				a, b = b, a
			}
			for b*2 < a {
				b *= 2
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
