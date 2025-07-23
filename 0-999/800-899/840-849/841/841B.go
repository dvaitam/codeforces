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
	odd := false
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x%2 != 0 {
			odd = true
		}
	}
	if odd {
		fmt.Println("First")
	} else {
		fmt.Println("Second")
	}
}
