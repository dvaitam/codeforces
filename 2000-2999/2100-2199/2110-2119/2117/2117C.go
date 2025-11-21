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
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		seen := make(map[int]struct{})
		ans := 0
		for _, x := range a {
			if x != a[0] && len(seen) == 0 {
				seen[a[0]] = struct{}{}
			}
			if _, ok := seen[x]; !ok {
				seen = map[int]struct{}{x: {}}
				ans++
			} else {
				seen[x] = struct{}{}
			}
		}
		ans++
		fmt.Fprintln(out, ans)
	}
}
