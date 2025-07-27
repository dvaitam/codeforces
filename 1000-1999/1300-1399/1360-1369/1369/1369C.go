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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		w := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &w[i])
		}

		sort.Ints(a)
		sort.Ints(w)

		l := 0
		r := n - 1
		var res int64

		// Assign largest numbers as maxima
		for i := 0; i < k; i++ {
			res += int64(a[r])
			if w[i] == 1 {
				res += int64(a[r])
			}
			r--
		}

		// Assign minima from smallest numbers
		for i := k - 1; i >= 0; i-- {
			if w[i] == 1 {
				continue
			}
			res += int64(a[l])
			l += w[i] - 1
		}

		fmt.Fprintln(writer, res)
	}
}
