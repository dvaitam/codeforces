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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	odd, even := 0, 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x%2 == 0 {
			even++
		} else {
			odd++
		}
	}

	var result int
	if even >= odd {
		result = odd
	} else {
		result = even + (odd-even)/3
	}
	fmt.Fprintln(writer, result)
}
