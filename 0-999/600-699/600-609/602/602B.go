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
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	const MaxVal = 100000
	count := make([]int, MaxVal+1)

	l := 0
	minVal := a[0]
	maxVal := a[0]
	count[a[0]] = 1
	best := 1

	for r := 1; r < n; r++ {
		v := a[r]
		count[v]++
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
		for maxVal-minVal > 1 {
			leftVal := a[l]
			count[leftVal]--
			if count[leftVal] == 0 {
				if leftVal == minVal {
					for minVal <= MaxVal && count[minVal] == 0 {
						minVal++
					}
				}
				if leftVal == maxVal {
					for maxVal >= 0 && count[maxVal] == 0 {
						maxVal--
					}
				}
			}
			l++
		}
		if r-l+1 > best {
			best = r - l + 1
		}
	}

	fmt.Fprintln(writer, best)

}
