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
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		var n int
		fmt.Fscan(reader, &n)
		if n == 3 {
			fmt.Fprintln(writer, "3 2 1")
			fmt.Fprintln(writer, "1 3 2")
			fmt.Fprintln(writer, "2 3 1")
		} else {
			v := make([]int, n)
			idx := make([]int, n)
			k := n
			for i := 0; i < n; i++ {
				v[i] = k
				k--
				idx[i] = i
			}
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					fmt.Fprint(writer, v[idx[j]])
					if j+1 < n {
						fmt.Fprint(writer, " ")
					}
					idx[j]--
					if idx[j] < 0 {
						idx[j] = n - 1
					}
				}
				fmt.Fprint(writer, "\n")
			}
		}
	}
}
