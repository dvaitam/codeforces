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
		freq := make(map[int]int64)
		for i := 1; i <= n; i++ {
			var x int
			fmt.Fscan(in, &x)
			key := x - i
			freq[key]++
		}
		var ans int64
		for _, v := range freq {
			ans += v * (v - 1) / 2
		}
		fmt.Fprintln(out, ans)
	}
}
