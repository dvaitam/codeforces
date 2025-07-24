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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	x := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &x[i])
	}

	diff := make([]int64, n+3)
	add := func(l, r int, v int64) {
		if l > r || r < 1 || l > n {
			return
		}
		if l < 1 {
			l = 1
		}
		if r > n {
			r = n
		}
		diff[l] += v
		diff[r+1] -= v
	}

	for i := 0; i < m-1; i++ {
		a := x[i]
		b := x[i+1]
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		D := int64(b - a)
		add(1, a-1, D)
		add(b+1, n, D)
		add(a+1, b-1, D-1)
		add(a, a, int64(b-1))
		add(b, b, int64(a))
	}

	cur := int64(0)
	for i := 1; i <= n; i++ {
		cur += diff[i]
		if i > 1 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, cur)
	}
	writer.WriteByte('\n')
}
