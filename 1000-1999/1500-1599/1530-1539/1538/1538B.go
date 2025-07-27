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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		sum := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
		}
		if sum%n != 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		avg := sum / n
		ans := 0
		for _, v := range a {
			if v > avg {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
