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
	// ans stores pairs of indices
	ans := make([][2]int, 0)
	// leftOver stores sizes of leftover groups
	leftOver := make([]int, 0)
	if n%2 == 1 {
		leftOver = append(leftOver, 1)
	}
	numPairs := n / 2
	gap := 1
	num := 1
	// First phase: pair by gaps
	for numPairs > 0 {
		changed := make([]bool, n+2)
		ptr := 1
		count := numPairs * num
		for i := 0; i < count; i++ {
			for changed[ptr] {
				ptr++
			}
			ans = append(ans, [2]int{ptr, ptr + gap})
			changed[ptr] = true
			changed[ptr+gap] = true
		}
		if numPairs%2 == 1 {
			leftOver = append(leftOver, num*2)
		}
		numPairs /= 2
		gap *= 2
		num *= 2
	}
	// remove last leftover as in C++ implementation
	if len(leftOver) > 0 {
		leftOver = leftOver[:len(leftOver)-1]
	}
	// if no further work, output
	if len(leftOver) < 2 {
		fmt.Fprintln(writer, len(ans))
		for _, p := range ans {
			fmt.Fprintln(writer, p[0], p[1])
		}
		return
	}
	// build remaining pairs
	first := leftOver[0]
	cur := make([]int, 0, first)
	for i := 0; i < first; i++ {
		cur = append(cur, n-i)
	}
	ptr := 1
	backPtr := n - first
	loPtr := 1
	for loPtr < len(leftOver) {
		size := leftOver[loPtr]
		newCur := make([]int, 0, size)
		if size == len(cur) {
			for _, z := range cur {
				ans = append(ans, [2]int{z, backPtr})
				newCur = append(newCur, backPtr)
				backPtr--
			}
			loPtr++
		} else {
			for _, z := range cur {
				ans = append(ans, [2]int{z, ptr})
				newCur = append(newCur, ptr)
				ptr++
			}
		}
		cur = append(cur, newCur...)
	}
	// output result
	fmt.Fprintln(writer, len(ans))
	for _, p := range ans {
		fmt.Fprintln(writer, p[0], p[1])
	}
}
