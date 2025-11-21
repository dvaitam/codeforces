package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func isPerfectSquare(x int) (int, bool) {
	r := int(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r, r*r == x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		val := 0
		for i := 0; i < len(s); i++ {
			val = val*10 + int(s[i]-'0')
		}
		r, ok := isPerfectSquare(val)
		if ok {
			fmt.Fprintln(out, 0, r)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
