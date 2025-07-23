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
	last := 0
	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t-last > 15 {
			fmt.Println(last + 15)
			return
		}
		last = t
	}
	if last+15 <= 90 {
		fmt.Println(last + 15)
	} else {
		fmt.Println(90)
	}
}
