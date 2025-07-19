package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 200001

var a [N]int
var req [N]int
var t [N]int

func add(x int) {
	for x < N {
		t[x]++
		x += x & -x
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		x, y := 0, 0
		for j := 17; j >= 0; j-- {
			nxt := x | (1 << j)
			if nxt < N && int64(a[i])*int64(nxt) <= int64(y+t[nxt]) {
				x = nxt
				y += t[nxt]
			}
		}
		x++
		add(x)
		req[i] = x
	}
	for ; q > 0; q-- {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		if y < req[x] {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
