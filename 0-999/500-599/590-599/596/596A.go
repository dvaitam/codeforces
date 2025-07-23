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
	xs := make([]int, n)
	ys := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &xs[i], &ys[i])
	}
	if n <= 1 {
		fmt.Println(-1)
		return
	}
	minX, maxX := xs[0], xs[0]
	minY, maxY := ys[0], ys[0]
	for i := 1; i < n; i++ {
		if xs[i] < minX {
			minX = xs[i]
		}
		if xs[i] > maxX {
			maxX = xs[i]
		}
		if ys[i] < minY {
			minY = ys[i]
		}
		if ys[i] > maxY {
			maxY = ys[i]
		}
	}
	width := maxX - minX
	height := maxY - minY
	if width == 0 || height == 0 {
		fmt.Println(-1)
	} else {
		fmt.Println(width * height)
	}
}
