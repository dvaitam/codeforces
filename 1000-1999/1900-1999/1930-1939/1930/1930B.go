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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		res := make([]int, 0, n)
		l, r := 1, n
		for l < r {
			res = append(res, l)
			res = append(res, r)
			l++
			r--
		}
		if l == r {
			res = append(res, l)
		}
		for i, v := range res {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, v)
		}
		out.WriteByte('\n')
	}
}
