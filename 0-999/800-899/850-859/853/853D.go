package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	costs := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &costs[i])
		total += costs[i]
	}

	var bonus int64
	var cash int64
	rem := total

	for _, c := range costs {
		if rem-bonus > c {
			cash += c
			bonus += c / 10
		} else {
			if bonus >= c {
				bonus -= c
			} else {
				cash += c - bonus
				bonus = 0
			}
		}
		rem -= c
	}

	fmt.Fprintln(out, cash)
}
