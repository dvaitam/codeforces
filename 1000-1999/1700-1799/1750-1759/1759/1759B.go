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
		var m, s int
		fmt.Fscan(reader, &m, &s)
		sumB := 0
		maxB := 0
		for i := 0; i < m; i++ {
			var x int
			fmt.Fscan(reader, &x)
			sumB += x
			if x > maxB {
				maxB = x
			}
		}
		total := sumB + s
		nSum := 0
		n := 0
		for nSum < total {
			n++
			nSum += n
		}
		if nSum == total && maxB <= n {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
