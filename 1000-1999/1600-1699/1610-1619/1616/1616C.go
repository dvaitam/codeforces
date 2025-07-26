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
		if n <= 2 {
			fmt.Fprintln(writer, 0)
			continue
		}
		best := 2
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				cnt := 2
				d1 := a[j] - a[i]
				d2 := j - i
				for k := 0; k < n; k++ {
					if k == i || k == j {
						continue
					}
					if (a[k]-a[i])*d2 == d1*(k-i) {
						cnt++
					}
				}
				if cnt > best {
					best = cnt
				}
			}
		}
		fmt.Fprintln(writer, n-best)
	}
}
