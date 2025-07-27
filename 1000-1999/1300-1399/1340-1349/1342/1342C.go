package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
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
		var a, b int64
		var q int
		fmt.Fscan(reader, &a, &b, &q)

		if a > b {
			a, b = b, a
		}
		lcm := a / gcd(a, b) * b
		period := int(lcm)
		pre := make([]int, period+1)
		for i := 0; i < period; i++ {
			if (int64(i)%a)%b != (int64(i)%b)%a {
				pre[i+1] = pre[i] + 1
			} else {
				pre[i+1] = pre[i]
			}
		}
		total := pre[period]

		calc := func(x int64) int64 {
			if x < 0 {
				return 0
			}
			q := x / lcm
			r := int(x % lcm)
			return int64(total)*q + int64(pre[r+1])
		}

		for i := 0; i < q; i++ {
			var l, r int64
			fmt.Fscan(reader, &l, &r)
			ans := calc(r) - calc(l-1)
			if i+1 == q {
				fmt.Fprint(writer, ans)
			} else {
				fmt.Fprint(writer, ans, " ")
			}
		}
		if t > 1 {
			fmt.Fprintln(writer)
		} else {
			fmt.Fprint(writer, "\n")
		}
	}
}
