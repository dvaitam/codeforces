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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		var y int64
		fmt.Fscan(in, &n, &y)
		c := make([]int, n)
		maxVal := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
			if c[i] > maxVal {
				maxVal = c[i]
			}
		}
		freq := make([]int, maxVal+2)
		for _, v := range c {
			freq[v]++
		}
		prefix := make([]int, maxVal+2)
		for i := 1; i <= maxVal; i++ {
			prefix[i] = prefix[i-1] + freq[i]
		}

		best := int64(0)
		hasBest := false
		totalN := int64(n)
		for x := 2; x <= maxVal+1; x++ {
			var sumCeil int64
			var sumMatch int64
			yVal := 1
			for {
				l := (yVal - 1) * x
				if l >= maxVal {
					break
				}
				r := yVal * x
				if r > maxVal {
					r = maxVal
				}
				count := prefix[r] - prefix[l]
				if count > 0 {
					sumCeil += int64(count * yVal)
					if yVal <= maxVal {
						if freq[yVal] < count {
							sumMatch += int64(freq[yVal])
						} else {
							sumMatch += int64(count)
						}
					}
				}
				yVal++
			}
			income := sumCeil - y*(totalN-sumMatch)
			if !hasBest || income > best {
				hasBest = true
				best = income
			}
		}

		fmt.Fprintln(out, best)
	}
}
