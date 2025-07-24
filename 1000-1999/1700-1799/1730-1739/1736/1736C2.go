package main

import (
	"bufio"
	"fmt"
	"os"
)

func goodCount(a []int) int64 {
	w := 1
	var ans int64
	for i := 0; i < len(a); i++ {
		val := i + 1 - a[i] + 1
		if val > w {
			w = val
		}
		if w < 1 {
			w = 1
		}
		ans += int64(i + 1 - w + 1)
	}
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var p, x int
		fmt.Fscan(reader, &p, &x)
		old := a[p-1]
		a[p-1] = x
		ans := goodCount(a)
		fmt.Fprintln(writer, ans)
		a[p-1] = old
	}
}
