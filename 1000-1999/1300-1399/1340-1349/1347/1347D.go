package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		l, r := 0, n-1
		moves := 0
		aliceTotal, bobTotal := 0, 0
		prev := 0
		aliceTurn := true
		for l <= r {
			cur := 0
			if aliceTurn {
				for l <= r && cur <= prev {
					cur += arr[l]
					l++
				}
				aliceTotal += cur
			} else {
				for l <= r && cur <= prev {
					cur += arr[r]
					r--
				}
				bobTotal += cur
			}
			prev = cur
			moves++
			aliceTurn = !aliceTurn
		}
		fmt.Fprintf(out, "%d %d %d\n", moves, aliceTotal, bobTotal)
	}
}
