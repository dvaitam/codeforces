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
		var n, k int
		fmt.Fscan(in, &n, &k)
		nums := make([]int, n)
		set := make(map[int]bool, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &nums[i])
			set[nums[i]] = true
		}
		found := false
		for _, v := range nums {
			if set[v+k] {
				found = true
				break
			}
		}
		if found {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
