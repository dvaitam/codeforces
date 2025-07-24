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
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			for x%2 == 0 {
				x /= 2
			}
			freq[x]++
		}

		bs := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &bs[i])
		}

		possible := true
		for _, b := range bs {
			for b%2 == 0 {
				b /= 2
			}
			for b > 0 && freq[b] == 0 {
				b /= 2
			}
			if b == 0 {
				possible = false
				break
			}
			freq[b]--
		}

		if possible {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
