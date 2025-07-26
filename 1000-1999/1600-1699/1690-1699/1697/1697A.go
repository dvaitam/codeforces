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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		sum := 0
		for i := 0; i < n; i++ {
			var a int
			fmt.Fscan(reader, &a)
			sum += a
		}
		need := sum - m
		if need < 0 {
			need = 0
		}
		fmt.Fprintln(writer, need)
	}
}
