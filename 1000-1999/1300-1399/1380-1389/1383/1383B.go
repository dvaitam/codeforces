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
		total := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			total ^= arr[i]
		}
		if total == 0 {
			fmt.Fprintln(out, "DRAW")
			continue
		}
		highest := 0
		for bit := 30; bit >= 0; bit-- {
			if (total>>bit)&1 == 1 {
				highest = bit
				break
			}
		}
		cnt1 := 0
		for _, v := range arr {
			if (v>>highest)&1 == 1 {
				cnt1++
			}
		}
		cnt0 := n - cnt1
		if cnt1%4 == 1 || (cnt1%4 == 3 && cnt0%2 == 1) {
			fmt.Fprintln(out, "WIN")
		} else {
			fmt.Fprintln(out, "LOSE")
		}
	}
}
