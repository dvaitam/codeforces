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
		sum := 0
		c1, c2 := 0, 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += x
			switch x % 3 {
			case 1:
				c1++
			case 2:
				c2++
			}
		}
		mod := sum % 3
		ans := 0
		if mod == 1 {
			if c1 >= 1 {
				ans = 1
			} else {
				ans = 2
			}
		} else if mod == 2 {
			if c2 >= 1 {
				ans = 1
			} else {
				ans = 2
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
