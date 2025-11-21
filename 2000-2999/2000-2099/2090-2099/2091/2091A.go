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
	required := [10]int{}
	for _, d := range []int{0, 1, 0, 3, 2, 0, 2, 5} {
		required[d]++
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		digits := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &digits[i])
		}
		cnt := [10]int{}
		ans := 0
	needMet:
		for i := 0; i < n; i++ {
			cnt[digits[i]]++
			ok := true
			for d := 0; d <= 9; d++ {
				if cnt[d] < required[d] {
					ok = false
					break
				}
			}
			if ok {
				ans = i + 1
				break needMet
			}
		}
		fmt.Fprintln(out, ans)
	}
}
