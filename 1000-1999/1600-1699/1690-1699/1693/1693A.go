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
		sum := 0
		zeroReached := false
		ok := true
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += x
			if sum < 0 {
				ok = false
			}
			if zeroReached && sum != 0 {
				ok = false
			}
			if sum == 0 {
				zeroReached = true
			}
		}
		if sum != 0 {
			ok = false
		}
		if ok {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
