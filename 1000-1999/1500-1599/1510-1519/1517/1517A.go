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
		if n%2050 != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		q := n / 2050
		sum := int64(0)
		for q > 0 {
			sum += q % 10
			q /= 10
		}
		fmt.Fprintln(writer, sum)
	}
}
