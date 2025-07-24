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

	var a, b, c int
	if _, err := fmt.Fscan(reader, &a, &b, &c); err != nil {
		return
	}
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	count := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x > b && x < c {
			count++
		}
	}
	fmt.Fprintln(writer, count)
}
