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
		var n int64
		fmt.Fscan(reader, &n)
		spf := int64(-1)
		for i := int64(2); i*i <= n; i++ {
			if n%i == 0 {
				spf = i
				break
			}
		}
		if spf == -1 {
			fmt.Fprintln(writer, 1, n-1)
		} else {
			a := n / spf
			fmt.Fprintln(writer, a, n-a)
		}
	}
}
