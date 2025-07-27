package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isPerfectSquare(x int) bool {
	r := int(math.Sqrt(float64(x)))
	return r*r == x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		ans := "NO"
		for i := 0; i < n; i++ {
			var val int
			fmt.Fscan(reader, &val)
			if !isPerfectSquare(val) {
				ans = "YES"
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
