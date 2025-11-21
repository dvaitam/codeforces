package main

import (
	"bufio"
	"fmt"
	"os"
)

func v2(x int64) int {
	cnt := 0
	for x%2 == 0 {
		cnt++
		x /= 2
	}
	return cnt
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var a, b int64
		fmt.Fscan(in, &a, &b)
		if b%2 == 1 {
			if a%2 == 0 {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, a*b+1)
			}
			continue
		}
		s := v2(b)
		if s == 1 {
			if a%2 == 1 {
				fmt.Fprintln(out, -1)
			} else {
				fmt.Fprintln(out, a*(b/2)+2)
			}
		} else {
			fmt.Fprintln(out, a*(b/2)+2)
		}
	}
}
