package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 1e9 + 7

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		sort.Ints(a)
		sort.Ints(b)

		res := int64(1)
		for i := n - 1; i >= 0; i-- {
			pos := sort.Search(len(a), func(idx int) bool { return a[idx] > b[i] })
			count := n - pos - (n - 1 - i)
			if count <= 0 {
				res = 0
				break
			}
			res = (res * int64(count)) % mod
		}

		fmt.Fprintln(writer, res)
	}
}
