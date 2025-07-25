package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		used := make([]bool, n-1)
		for x := 1; x <= n; x++ {
			// find position of x
			pos := 0
			for pos < n && a[pos] != x {
				pos++
			}
			// move x left while beneficial and operation not used
			for pos > 0 && !used[pos-1] && a[pos-1] > a[pos] {
				a[pos], a[pos-1] = a[pos-1], a[pos]
				used[pos-1] = true
				pos--
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, a[i])
		}
		fmt.Fprintln(writer)
	}
}
