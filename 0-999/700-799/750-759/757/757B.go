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
	const maxVal = 100000
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	// Precompute smallest prime factors up to maxVal
	spf := make([]int, maxVal+1)
	for i := 2; i <= maxVal; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxVal; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}

	count := make(map[int]int)
	best := 0
	for _, x := range arr {
		if x <= 1 {
			continue
		}
		seen := make(map[int]bool)
		for x > 1 {
			p := spf[x]
			if p == 0 {
				p = x
			}
			if !seen[p] {
				count[p]++
				if count[p] > best {
					best = count[p]
				}
				seen[p] = true
			}
			for x%p == 0 {
				x /= p
			}
		}
	}

	if best == 0 {
		fmt.Fprintln(writer, 1)
	} else {
		fmt.Fprintln(writer, best)
	}
}
