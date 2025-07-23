package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, h int
	if _, err := fmt.Fscan(reader, &n, &h); err != nil {
		return
	}
	width := 0
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(reader, &a)
		if a > h {
			width += 2
		} else {
			width++
		}
	}
	fmt.Println(width)
}
