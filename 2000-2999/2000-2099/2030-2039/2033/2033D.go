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

		sum := int64(0)
		ans := 0
		seen := make(map[int64]struct{})
		seen[0] = struct{}{}

		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			sum += x
			if _, ok := seen[sum]; ok {
				ans++
				seen = make(map[int64]struct{})
			}
			seen[sum] = struct{}{}
		}

		fmt.Fprintln(out, ans)
	}
}
