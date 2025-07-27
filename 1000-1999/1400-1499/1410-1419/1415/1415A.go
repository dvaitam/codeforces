package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for i := 0; i < t; i++ {
		var n, m, r, c int64
		fmt.Fscan(reader, &n, &m, &r, &c)
		d1 := abs(1-r) + abs(1-c)
		d2 := abs(1-r) + abs(m-c)
		d3 := abs(n-r) + abs(1-c)
		d4 := abs(n-r) + abs(m-c)
		ans := d1
		if d2 > ans {
			ans = d2
		}
		if d3 > ans {
			ans = d3
		}
		if d4 > ans {
			ans = d4
		}
		fmt.Fprintln(writer, ans)
	}
}
