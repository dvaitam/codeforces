package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var g, c, l int
	if _, err := fmt.Fscan(in, &g, &c, &l); err != nil {
		return
	}

	minVal := g
	maxVal := g
	if c < minVal {
		minVal = c
	}
	if l < minVal {
		minVal = l
	}
	if c > maxVal {
		maxVal = c
	}
	if l > maxVal {
		maxVal = l
	}

	if maxVal-minVal >= 10 {
		fmt.Println("check again")
		return
	}

	sum := g + c + l
	median := sum - minVal - maxVal
	fmt.Printf("final %d\n", median)
}
