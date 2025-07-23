package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	// Maximum exponent from input determines array size.
	maxW := 0
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &weights[i])
		if weights[i] > maxW {
			maxW = weights[i]
		}
	}

	// Extra space for carries: log2(n) < 20, so add 70 to be safe.
	cnt := make([]int, maxW+70)
	for _, w := range weights {
		cnt[w]++
	}

	for i := 0; i < len(cnt)-1; i++ {
		pairs := cnt[i] / 2
		if pairs > 0 {
			cnt[i+1] += pairs
			cnt[i] %= 2
		}
	}

	steps := 0
	for _, c := range cnt {
		steps += c
	}
	fmt.Fprintln(writer, steps)
}
