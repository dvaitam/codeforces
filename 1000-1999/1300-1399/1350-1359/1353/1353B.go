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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		sort.Ints(a)
		sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
		for i := 0; i < k && i < n; i++ {
			if a[i] < b[i] {
				a[i], b[i] = b[i], a[i]
			} else {
				break
			}
		}
		sum := 0
		for _, v := range a {
			sum += v
		}
		fmt.Fprintln(writer, sum)
	}
}
