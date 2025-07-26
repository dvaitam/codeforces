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
		count := make([]int, 10)
		digits := make([]int, 0, 30)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			d := x % 10
			if count[d] < 3 {
				count[d]++
				digits = append(digits, d)
			}
		}
		found := false
		for i := 0; i < len(digits) && !found; i++ {
			for j := i + 1; j < len(digits) && !found; j++ {
				for k := j + 1; k < len(digits) && !found; k++ {
					if (digits[i]+digits[j]+digits[k])%10 == 3 {
						found = true
					}
				}
			}
		}
		if found {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
