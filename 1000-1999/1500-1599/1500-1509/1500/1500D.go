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
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	c := make([][]int, n)
	for i := 0; i < n; i++ {
		c[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &c[i][j])
		}
	}

	res := make([]int, n)
	for k := 1; k <= n; k++ {
		count := 0
		for i := 0; i <= n-k; i++ {
			for j := 0; j <= n-k; j++ {
				colors := make(map[int]struct{})
				ok := true
				for x := 0; x < k && ok; x++ {
					for y := 0; y < k && ok; y++ {
						colors[c[i+x][j+y]] = struct{}{}
						if len(colors) > q {
							ok = false
						}
					}
				}
				if ok {
					count++
				}
			}
		}
		res[k-1] = count
	}

	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
