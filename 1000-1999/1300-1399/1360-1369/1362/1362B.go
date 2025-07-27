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
		arr := make([]int, n)
		set := make(map[int]bool, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			set[arr[i]] = true
		}
		ans := -1
		for k := 1; k < 1024; k++ {
			ok := true
			for _, v := range arr {
				if !set[v^k] {
					ok = false
					break
				}
			}
			if ok {
				ans = k
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}
