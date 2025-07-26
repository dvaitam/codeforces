package main

import (
	"bufio"
	"fmt"
	"os"
)

func minCost(s string) int {
	n := len(s)
	// find first '1'
	first := -1
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			first = i
			break
		}
	}
	if first == -1 {
		// all zeros
		return 0
	}
	last := -1
	for i := n - 1; i >= 0; i-- {
		if s[i] == '0' {
			last = i
			break
		}
	}
	if last == -1 || first > last {
		// already sorted
		return 0
	}
	zerosBefore := first
	onesAfter := n - 1 - last
	zerosInside := 0
	onesInside := 0
	for i := first; i <= last; i++ {
		if s[i] == '0' {
			zerosInside++
		} else {
			onesInside++
		}
	}
	delta := onesInside - zerosInside
	diff := 0
	if delta > zerosBefore {
		diff = delta - zerosBefore
	} else if -delta > onesAfter {
		diff = -delta - onesAfter
	}
	return diff + 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)
		fmt.Fprintln(writer, minCost(s))
	}
}
