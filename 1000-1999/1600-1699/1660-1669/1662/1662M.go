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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		maxR, maxW := 0, 0
		pairs := make([][2]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &pairs[i][0], &pairs[i][1])
			if pairs[i][0] > maxR {
				maxR = pairs[i][0]
			}
			if pairs[i][1] > maxW {
				maxW = pairs[i][1]
			}
		}

		if maxR+maxW > n {
			fmt.Fprintln(out, "IMPOSSIBLE")
			continue
		}
		res := make([]byte, n)
		for i := 0; i < maxR; i++ {
			res[i] = 'R'
		}
		for i := maxR; i < n; i++ {
			res[i] = 'W'
		}
		fmt.Fprintln(out, string(res))
	}
}
