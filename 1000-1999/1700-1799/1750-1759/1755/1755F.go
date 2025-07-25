package main

import (
	"bufio"
	"fmt"
	"os"
)

func divisors(x int) []int {
	divs := []int{}
	for i := 1; i*i <= x; i++ {
		if x%i == 0 {
			divs = append(divs, i)
			if i*i != x {
				divs = append(divs, x/i)
			}
		}
	}
	return divs
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			freq[arr[i]]++
		}

		nHalf := n / 2
		impossible := false
		for _, c := range freq {
			if c >= nHalf {
				impossible = true
				break
			}
		}
		if impossible {
			fmt.Fprintln(writer, -1)
			continue
		}

		candidates := make(map[int]struct{})
		maxK := 1
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				diff := arr[i] - arr[j]
				if diff < 0 {
					diff = -diff
				}
				if diff == 0 {
					continue
				}
				for _, d := range divisors(diff) {
					if d <= maxK {
						continue
					}
					if _, ok := candidates[d]; ok {
						continue
					}
					candidates[d] = struct{}{}
					counts := make(map[int]int)
					valid := false
					for _, val := range arr {
						r := val % d
						if r < 0 {
							r += d
						}
						counts[r]++
						if counts[r] >= nHalf {
							valid = true
							break
						}
					}
					if valid && d > maxK {
						maxK = d
					}
				}
			}
		}
		fmt.Fprintln(writer, maxK)
	}
}
