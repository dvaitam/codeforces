package main

import (
	"bufio"
	"fmt"
	"os"
)

func norm(x, mod int) int {
	x %= mod
	if x < 0 {
		x += mod
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	odd := 0
	even := 1 % n // in case n==1? but n>=2
	for i := 0; i < q; i++ {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var x int
			fmt.Fscan(reader, &x)
			x = norm(x, n)
			if x%2 == 0 {
				odd = norm(odd-x, n)
				even = norm(even-x, n)
			} else {
				newOdd := norm(even-(x+1), n)
				newEven := norm(odd-(x-1), n)
				odd, even = newOdd, newEven
			}
		} else {
			odd, even = even, odd
		}
	}
	res := make([]int, n)
	for i := 0; i < n; i += 2 {
		res[i] = norm(odd+i, n) + 1
		if i+1 < n {
			res[i+1] = norm(even+i, n) + 1
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, res[i])
	}
	writer.WriteByte('\n')
}
