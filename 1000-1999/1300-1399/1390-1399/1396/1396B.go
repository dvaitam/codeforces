package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(reader, &n)
		sum := 0
		maxv := 0
		for j := 0; j < n; j++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += x
			if x > maxv {
				maxv = x
			}
		}
		// If the largest pile has more stones than all others combined, T wins immediately.
		// Otherwise, if total stones is odd, T wins; else HL wins.
		if maxv > sum-maxv || sum%2 == 1 {
			fmt.Println("T")
		} else {
			fmt.Println("HL")
		}
	}
}
