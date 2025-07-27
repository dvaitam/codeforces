package main

import (
	"bufio"
	"fmt"
	"os"
)

func waitTime(p, x int64) int64 {
	if p%x == 0 {
		return 0
	}
	return x - p%x
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
		var p, a, b, c int64
		if _, err := fmt.Fscan(reader, &p, &a, &b, &c); err != nil {
			return
		}
		wa := waitTime(p, a)
		wb := waitTime(p, b)
		wc := waitTime(p, c)
		ans := wa
		if wb < ans {
			ans = wb
		}
		if wc < ans {
			ans = wc
		}
		fmt.Fprintln(writer, ans)
	}
}
