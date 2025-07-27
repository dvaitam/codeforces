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
		for k := int64(2); ; k++ {
			denom := (int64(1) << k) - 1
			if n%denom == 0 {
				fmt.Fprintln(writer, n/denom)
				break
			}
		}
	}
}
