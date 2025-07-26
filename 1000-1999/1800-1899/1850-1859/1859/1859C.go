package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxCost(n int) int {
	best := 0
	for r := 0; r < n; r++ {
		sum := 0
		maxProd := 0
		for i := 1; i <= r; i++ {
			prod := i * i
			sum += prod
			if prod > maxProd {
				maxProd = prod
			}
		}
		for j := 1; j <= n-r; j++ {
			val := n - j + 1
			idx := r + j
			prod := idx * val
			sum += prod
			if prod > maxProd {
				maxProd = prod
			}
		}
		cost := sum - maxProd
		if cost > best {
			best = cost
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t, n int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		fmt.Fscan(reader, &n)
		fmt.Fprintln(writer, maxCost(n))
	}
}
