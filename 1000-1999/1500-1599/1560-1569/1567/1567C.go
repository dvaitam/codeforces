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
		var a, b int64
		multA, multB := int64(1), int64(1)
		pos := 0
		for n > 0 {
			digit := n % 10
			if pos%2 == 0 {
				a += digit * multA
				multA *= 10
			} else {
				b += digit * multB
				multB *= 10
			}
			pos++
			n /= 10
		}
		ans := (a+1)*(b+1) - 2
		fmt.Fprintln(writer, ans)
	}
}
