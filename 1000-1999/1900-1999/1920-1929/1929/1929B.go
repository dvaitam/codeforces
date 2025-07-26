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
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		if k == 4*n-2 {
			fmt.Fprintln(writer, 2*n)
		} else {
			fmt.Fprintln(writer, (k+1)/2)
		}
	}
}
