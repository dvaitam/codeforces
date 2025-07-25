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
		freq := make([]int, 101)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			if x >= 0 && x < len(freq) {
				freq[x]++
			}
		}
		ans := 0
		for _, c := range freq {
			ans += c / 3
		}
		fmt.Fprintln(out, ans)
	}
}
