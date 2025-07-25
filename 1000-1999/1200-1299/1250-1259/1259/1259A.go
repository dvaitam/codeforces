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
		var n int
		fmt.Fscan(reader, &n)

		digits := 0
		temp := n
		for temp > 0 {
			digits++
			temp /= 10
		}

		ones := 0
		for i := 0; i < digits; i++ {
			ones = ones*10 + 1
		}

		ans := (digits-1)*9 + n/ones
		fmt.Fprintln(writer, ans)
	}
}
