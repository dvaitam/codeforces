package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var N, M int
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return
	}
	bestK := 1
	bestTotal := -1
	for K := 1; K <= M; K++ {
		boxes := make(map[int]bool)
		total := 0
		for j := 1; j <= N; j++ {
			boxes[j] = true
			y := (j*K-1)%N + 1
			if !boxes[y] {
				total += y
				boxes[y] = true
			}
		}
		if total > bestTotal {
			bestTotal = total
			bestK = K
		}
	}
	fmt.Println(bestK)
}
