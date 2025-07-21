package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var x, y int
		fmt.Fscan(in, &x, &y)
		if x < y {
			fmt.Printf("%d %d\n", x, y)
		} else {
			fmt.Printf("%d %d\n", y, x)
		}
	}
}
