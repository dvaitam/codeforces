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
		var k int64
		fmt.Fscan(reader, &n, &k)
		rem := make([]int64, n)
		for i := 0; i < n; i++ {
			var a int64
			fmt.Fscan(reader, &a)
			rem[i] = a % k
		}
		dp0, dp1 := rem[0], rem[0]+k
		val0, val1 := rem[0], rem[0]+k
		for i := 1; i < n; i++ {
			cur := rem[i]
			c0 := dp0
			if cur > val0 {
				c0 += cur - val0
			}
			c1 := dp1
			if cur > val1 {
				c1 += cur - val1
			}
			ndp0 := c0
			if c1 < ndp0 {
				ndp0 = c1
			}

			up := cur + k
			c0 = dp0
			if up > val0 {
				c0 += up - val0
			}
			c1 = dp1
			if up > val1 {
				c1 += up - val1
			}
			ndp1 := c0
			if c1 < ndp1 {
				ndp1 = c1
			}
			dp0, dp1 = ndp0, ndp1
			val0, val1 = cur, up
		}
		if dp1 < dp0 {
			fmt.Fprintln(writer, dp1)
		} else {
			fmt.Fprintln(writer, dp0)
		}
	}
}
