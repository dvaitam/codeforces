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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		cnt := [3]int{}
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
			cnt[x%3]++
		}
		target := n / 3
		moves := 0
		for i := 0; i < 6; i++ { // enough iterations to balance all
			idx := i % 3
			if cnt[idx] > target {
				diff := cnt[idx] - target
				cnt[idx] -= diff
				cnt[(idx+1)%3] += diff
				moves += diff
			}
		}
		fmt.Fprintln(writer, moves)
	}
}
