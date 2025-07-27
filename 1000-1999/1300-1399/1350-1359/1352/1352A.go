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
		res := make([]int, 0)
		base := 1
		for n > 0 {
			digit := n % 10
			if digit != 0 {
				res = append(res, digit*base)
			}
			n /= 10
			base *= 10
		}
		fmt.Fprintln(out, len(res))
		for i := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, res[i])
		}
		fmt.Fprintln(out)
	}
}
