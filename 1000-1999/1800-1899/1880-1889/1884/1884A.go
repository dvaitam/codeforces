package main

import (
	"bufio"
	"fmt"
	"os"
)

func digitSum(x int) int {
	sum := 0
	for x > 0 {
		sum += x % 10
		x /= 10
	}
	return sum
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x, k int
		fmt.Fscan(reader, &x, &k)
		for {
			if digitSum(x)%k == 0 {
				fmt.Fprintln(writer, x)
				break
			}
			x++
		}
	}
}
