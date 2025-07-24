package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func countLuxury(n int64) int64 {
	if n <= 0 {
		return 0
	}
	m := int64(math.Sqrt(float64(n)))
	for (m+1)*(m+1) <= n { // ensure floor sqrt
		m++
	}
	for m*m > n {
		m--
	}
	ans := (m - 1) * 3
	if m*m <= n {
		ans++
	}
	if m*m+m <= n {
		ans++
	}
	if m*m+2*m <= n {
		ans++
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		fmt.Fprintln(out, countLuxury(r)-countLuxury(l-1))
	}
}
