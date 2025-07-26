package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var r int
	if _, err := fmt.Fscan(in, &r); err != nil {
		return
	}
	thresholds := []int{-1000, 1200, 1400, 1600, 1900, 2100, 2300, 2400, 2600, 3000}
	for p := len(thresholds) - 1; p > 0; p-- {
		if r >= thresholds[p-1] {
			fmt.Println(thresholds[p])
			return
		}
	}
}
