package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	type frac struct {
		num, den int
	}
	arr := make([]frac, m)
	count := make(map[frac]int)
	for i := 0; i < m; i++ {
		var expr string
		fmt.Fscan(reader, &expr)
		var a, b, c int
		fmt.Sscanf(expr, "(%d+%d)/%d", &a, &b, &c)
		n := a + b
		g := gcd(n, c)
		f := frac{num: n / g, den: c / g}
		arr[i] = f
		count[f]++
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, count[arr[i]])
	}
	writer.WriteByte('\n')
}
