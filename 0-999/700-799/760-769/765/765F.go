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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	var m int
	fmt.Fscan(reader, &m)
	for ; m > 0; m-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		b := make([]int, r-l+1)
		copy(b, arr[l:r+1])
		sort.Ints(b)
		minDiff := b[1] - b[0]
		for i := 2; i < len(b); i++ {
			d := b[i] - b[i-1]
			if d < minDiff {
				minDiff = d
			}
		}
		fmt.Fprintln(writer, minDiff)
	}
}
