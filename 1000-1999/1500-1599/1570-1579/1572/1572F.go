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

	var n, q int
	fmt.Fscan(in, &n, &q)

	h := make([]int, n+1)
	w := make([]int, n+1)
	for i := 1; i <= n; i++ {
		h[i] = 0
		w[i] = i
	}

	maxH := 0

	for ; q > 0; q-- {
		var p int
		fmt.Fscan(in, &p)
		if p == 1 {
			var c, g int
			fmt.Fscan(in, &c, &g)
			maxH++
			h[c] = maxH
			w[c] = g
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			b := make([]int, n+1)
			for j := 1; j <= n; j++ {
				count := 0
				for i := 1; i <= j; i++ {
					if j <= w[i] {
						valid := true
						hi := h[i]
						for k := i + 1; k <= j; k++ {
							if h[k] >= hi {
								valid = false
								break
							}
						}
						if valid {
							count++
						}
					}
				}
				b[j] = count
			}
			sum := 0
			for j := l; j <= r; j++ {
				sum += b[j]
			}
			fmt.Fprintln(out, sum)
		}
	}
}
