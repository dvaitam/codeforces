package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	for i := 0; i < m; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
	}
	a := 1
	for i := 0; i < n; i++ {
		out.WriteByte(byte('0' + a))
		a = 1 - a
	}
}
