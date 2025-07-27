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
		var n int
		fmt.Fscan(reader, &n)
		count := 0
		for length := 1; length <= 9; length++ {
			for digit := 1; digit <= 9; digit++ {
				val := 0
				for i := 0; i < length; i++ {
					val = val*10 + digit
				}
				if val <= n {
					count++
				}
			}
		}
		fmt.Fprintln(writer, count)
	}
}
