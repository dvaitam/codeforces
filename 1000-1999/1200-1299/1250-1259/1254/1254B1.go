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
	a := make([]int, n)
	positions := make([]int, 0)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		if a[i] == 1 {
			positions = append(positions, i+1)
		}
	}

	total := len(positions)
	if total <= 1 {
		fmt.Fprintln(writer, -1)
		return
	}

	// get prime factors of total
	factors := make([]int, 0)
	x := total
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			factors = append(factors, p)
			for x%p == 0 {
				x /= p
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}

	var best int64 = -1
	for _, k := range factors {
		var cost int64
		for i := 0; i < total; i += k {
			median := positions[i+k/2]
			for j := i; j < i+k; j++ {
				if positions[j] > median {
					cost += int64(positions[j] - median)
				} else {
					cost += int64(median - positions[j])
				}
			}
		}
		if best == -1 || cost < best {
			best = cost
		}
	}

	fmt.Fprintln(writer, best)
}
