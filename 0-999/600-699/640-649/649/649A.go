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
	cnt := 0
	maxPow := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		// compute largest power of two dividing x
		pow := x & -x
		if pow > maxPow {
			maxPow = pow
			cnt = 1
		} else if pow == maxPow {
			cnt++
		}
	}
	fmt.Printf("%d %d\n", maxPow, cnt)
}
