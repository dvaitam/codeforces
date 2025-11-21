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
		var n int64
		fmt.Fscan(in, &n)
		ans := make([]int64, 0)
		pow := int64(10)
		for k := 1; k <= 18; k++ {
			denom := pow + 1
			if n%denom == 0 {
				x := n / denom
				if x > 0 {
					ans = append(ans, x)
				}
			}
			pow *= 10
		}
		if len(ans) == 0 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintln(out, len(ans))
			for i, x := range ans {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, x)
			}
			fmt.Fprintln(out)
		}
	}
}
