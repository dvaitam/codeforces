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

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	a := make([][]int64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}

	var s int64 = 1
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			s = lcm(s, a[i][j])
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			var val int64
			if (i&1)^(j&1) == 1 {
				x := a[i][j]
				val = x*x*x*x + s
			} else {
				val = s
			}
			if j > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, val)
		}
		writer.WriteByte('\n')
	}
}
