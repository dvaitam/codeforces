package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var x, y int64
	fmt.Fscan(in, &n, &x, &y)
	var s string
	fmt.Fscan(in, &s)

	zeros := 0
	i := 0
	for i < n {
		if s[i] == '0' {
			zeros++
			for i < n && s[i] == '0' {
				i++
			}
		} else {
			i++
		}
	}
	if zeros == 0 {
		fmt.Println(0)
		return
	}
	cost1 := int64(zeros) * y
	cost2 := y + int64(zeros-1)*x
	if cost1 < cost2 {
		fmt.Println(cost1)
	} else {
		fmt.Println(cost2)
	}
}
