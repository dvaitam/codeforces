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
	seenX := make([]bool, 101)
	seenY := make([]bool, 101)
	cntX, cntY := 0, 0
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		if x >= 0 && x <= 100 && !seenX[x] {
			seenX[x] = true
			cntX++
		}
		if y >= 0 && y <= 100 && !seenY[y] {
			seenY[y] = true
			cntY++
		}
	}
	if cntX < cntY {
		fmt.Println(cntX)
	} else {
		fmt.Println(cntY)
	}
}
