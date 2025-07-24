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
		sum := 0
		odd := 0
		even := 0
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += x
			if x%2 == 0 {
				even++
			} else {
				odd++
			}
		}
		if sum%2 == 1 || (odd > 0 && even > 0) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
