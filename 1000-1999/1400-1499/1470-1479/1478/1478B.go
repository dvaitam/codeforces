package main

import (
	"bufio"
	"fmt"
	"os"
)

func containsDigit(x, d int) bool {
	for x > 0 {
		if x%10 == d {
			return true
		}
		x /= 10
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var q, d int
		fmt.Fscan(reader, &q, &d)
		for i := 0; i < q; i++ {
			var a int
			fmt.Fscan(reader, &a)
			ok := false
			if a >= d*10 {
				ok = true
			} else {
				for k := 0; k <= 10; k++ {
					v := a - k*d
					if v < 0 {
						break
					}
					if containsDigit(v, d) {
						ok = true
						break
					}
				}
			}
			if ok {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
