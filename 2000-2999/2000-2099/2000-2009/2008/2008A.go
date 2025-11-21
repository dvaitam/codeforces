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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b int
		fmt.Fscan(in, &a, &b)
		ans := "No"
		if a%2 == 0 {
			target := a/2 + b
			for y := 0; y <= b; y++ {
				x := target - 2*y
				if x >= 0 && x <= a {
					ans = "Yes"
					break
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
