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

	if n <= 0 {
		fmt.Println(0)
		return
	}

	var x int64
	fmt.Fscan(in, &x)
	minVal := x

	for i := 1; i < n; i++ {
		fmt.Fscan(in, &x)
		if x < minVal {
			minVal = x
		}
	}

	fmt.Println(minVal)
}
