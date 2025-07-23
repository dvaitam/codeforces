package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, bx int
	if _, err := fmt.Fscan(in, &n, &bx); err != nil {
		return
	}
	var valX int64
	for i := 0; i < n; i++ {
		var d int64
		fmt.Fscan(in, &d)
		valX = valX*int64(bx) + d
	}
	var m, by int
	if _, err := fmt.Fscan(in, &m, &by); err != nil {
		return
	}
	var valY int64
	for i := 0; i < m; i++ {
		var d int64
		fmt.Fscan(in, &d)
		valY = valY*int64(by) + d
	}
	if valX < valY {
		fmt.Println("<")
	} else if valX > valY {
		fmt.Println(">")
	} else {
		fmt.Println("=")
	}
}
