package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		cnt := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			cnt[x]++
		}

		ans := 0
		for x, f := range cnt {
			if x == 0 {
				ans += f
			} else {
				if f < x {
					ans += f
				} else {
					ans += f - x
				}
			}
		}

		fmt.Println(ans)
	}
}

