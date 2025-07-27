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
	color := "blue"
	locked := false
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		switch s {
		case "lock":
			locked = true
		case "unlock":
			locked = false
		default:
			if !locked {
				color = s
			}
		}
	}
	fmt.Println(color)
}
