package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var sum int64
	for i := 0; i < n; i++ {
		var t int64
		fmt.Fscan(in, &t)
		sum += t
	}
	var m int
	fmt.Fscan(in, &m)
	for i := 0; i < m; i++ {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		if sum <= r {
			if sum < l {
				fmt.Println(l)
			} else {
				fmt.Println(sum)
			}
			return
		}
	}
	fmt.Println(-1)
}
