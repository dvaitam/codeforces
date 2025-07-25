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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k, q int
		fmt.Fscan(reader, &n, &k, &q)
		a := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &b[i])
		}
		for ; q > 0; q-- {
			var d int
			fmt.Fscan(reader, &d)
			if d == 0 {
				fmt.Fprintln(writer, 0)
				continue
			}
			idx := sort.SearchInts(a, d)
			if idx == len(a) {
				fmt.Fprintln(writer, b[k-1])
				continue
			}
			prevA, prevB := 0, 0
			if idx > 0 {
				prevA = a[idx-1]
				prevB = b[idx-1]
			}
			num := (d - prevA) * (b[idx] - prevB)
			den := a[idx] - prevA
			ans := prevB + num/den
			fmt.Fprintln(writer, ans)
		}
	}
}
