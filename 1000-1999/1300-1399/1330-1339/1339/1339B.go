package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Ints(a)
		// build result: start from middle, then take elements in alternating left/right jumps
		res := make([]int, 0, n)
		idx := n / 2
		for j := 1; j <= n; j++ {
			res = append(res, a[idx])
			if j%2 == 1 {
				idx -= j
			} else {
				idx += j
			}
		}
		// output
		for i, v := range res {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, v)
		}
		writer.WriteByte('\n')
	}
}
