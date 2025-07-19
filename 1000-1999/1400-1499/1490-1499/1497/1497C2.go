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

	var t, n, k int64
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		fmt.Fscan(reader, &n, &k)
		d := k - 3
		n2 := n - d
		// Determine three base parts
		if n2%3 == 0 {
			a := n2 / 3
			fmt.Fprint(writer, a, " ", a, " ", a, " ")
		} else if n2%2 == 0 {
			half := n2 / 2
			if half%2 == 1 {
				x := (n2 - 2) / 2
				fmt.Fprint(writer, x, " ", x, " ", 2, " ")
			} else {
				fmt.Fprint(writer, half, " ", half/2, " ", half/2, " ")
			}
		} else {
			x := (n2 - 1) / 2
			fmt.Fprint(writer, x, " ", x, " ", 1, " ")
		}
		// Append ones for remaining parts
		for i := int64(0); i < d; i++ {
			fmt.Fprint(writer, 1, " ")
		}
		fmt.Fprint(writer, "\n")
	}
}
