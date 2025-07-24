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

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n, r int
		fmt.Fscan(reader, &n, &r)
		x := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &x[i])
		}
		sort.Slice(x, func(i, j int) bool { return x[i] > x[j] })
		uniq := []int{x[0]}
		for i := 1; i < n; i++ {
			if x[i] != x[i-1] {
				uniq = append(uniq, x[i])
			}
		}
		shots := 0
		shift := 0
		for _, pos := range uniq {
			if pos-shift <= 0 {
				break
			}
			shots++
			shift += r
		}
		fmt.Fprintln(writer, shots)
	}
}
