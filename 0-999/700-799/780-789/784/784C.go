package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	maxVal := 0
	last := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		if x > maxVal {
			maxVal = x
		}
		last = x
	}
	fmt.Println(maxVal ^ last)
}
