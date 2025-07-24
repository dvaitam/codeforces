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
		var product int64 = 1
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			product *= x
		}
		ans := product + int64(n-1)
		fmt.Fprintln(out, ans*2022)
	}
}
