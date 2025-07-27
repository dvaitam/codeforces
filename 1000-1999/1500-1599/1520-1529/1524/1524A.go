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

	var g int
	if _, err := fmt.Fscan(in, &g); err != nil {
		return
	}
	for ; g > 0; g-- {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return
		}
		weights := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &weights[i])
		}
		var k int
		fmt.Fscan(in, &k)
		marked := make([]bool, n)
		for i := 0; i < k; i++ {
			var v int
			fmt.Fscan(in, &v)
			if v >= 1 && v <= n {
				marked[v-1] = true
			}
		}
		var m int
		fmt.Fscan(in, &m)
		for i := 0; i < m; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			// edges ignored in this simple solution
			_ = x
			_ = y
		}

		additional := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if !marked[i] {
				additional = append(additional, i+1)
			}
		}

		fmt.Fprint(out, len(additional))
		for _, v := range additional {
			fmt.Fprintf(out, " %d", v)
		}
		fmt.Fprintln(out)

		fmt.Fprintln(out, n)
		for i := 1; i <= n; i++ {
			fmt.Fprintf(out, "1 %d\n", i)
		}
	}
}
