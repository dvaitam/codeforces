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

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	first := make([]int, 51)
	for i := 1; i <= 50; i++ {
		first[i] = -1
	}
	for i := 1; i <= n; i++ {
		var c int
		fmt.Fscan(reader, &c)
		if first[c] == -1 {
			first[c] = i
		}
	}

	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(reader, &t)
		pos := first[t]
		fmt.Fprint(writer, pos)
		if i+1 < q {
			fmt.Fprint(writer, " ")
		}
		for c := 1; c <= 50; c++ {
			if first[c] != -1 && first[c] < pos {
				first[c]++
			}
		}
		first[t] = 1
	}
	fmt.Fprintln(writer)
}
