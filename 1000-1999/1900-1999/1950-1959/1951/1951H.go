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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var k int
		fmt.Fscan(in, &k)
		n := 1 << k
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		for t := 1; t <= k; t++ {
			ans := n - (1 << t) + 1
			if t > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans)
		}
		fmt.Fprintln(out)
	}
}
