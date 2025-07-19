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
		var n, k int
		fmt.Fscan(in, &n, &k)
		ans := 0
		if k >= n {
			k -= n
			ans++
		}
		for i := n - 1; i >= 1; i-- {
			if k >= i {
				k -= i
				ans++
			}
			if k >= i {
				k -= i
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
