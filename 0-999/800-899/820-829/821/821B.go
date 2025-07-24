package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var m, b int64
	fmt.Fscan(in, &m, &b)

	var best int64
	for y := int64(0); y <= b; y++ {
		x := m * (b - y)
		sumX := x * (x + 1) / 2
		sumY := y * (y + 1) / 2
		bananas := (y+1)*sumX + (x+1)*sumY
		if bananas > best {
			best = bananas
		}
	}

	fmt.Println(best)
}
