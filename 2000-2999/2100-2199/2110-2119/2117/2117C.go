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
		fmt.Fscan(in, &n)
		a := make([]int, n)
		maxVal := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}

		remain := make([]int, maxVal+2)
		for _, x := range a {
			remain[x]++
		}

		firstSeg := make([]int, maxVal+2)
		seenSeg := make([]int, maxVal+2)
		for i := range firstSeg {
			firstSeg[i] = -1
			seenSeg[i] = -1
		}

		distinctSeen := 0
		currSeg := 0
		needLeft := 0
		bad := 0
		cuts := 0

		for i := 0; i < n; i++ {
			x := a[i]
			remain[x]--
			if remain[x] == 0 {
				bad++
			}

			if firstSeg[x] == -1 {
				firstSeg[x] = currSeg
				distinctSeen++
			}

			if seenSeg[x] != currSeg {
				seenSeg[x] = currSeg
				if firstSeg[x] < currSeg {
					needLeft--
				}
			}

			if needLeft == 0 && bad == 0 && i < n-1 {
				cuts++
				currSeg++
				needLeft = distinctSeen
			}
		}

		fmt.Fprintln(out, cuts+1)
	}
}
