package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	votes := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &votes[i])
	}
	count := 0
	for {
		maxIdx := 0
		for i := 1; i < n; i++ {
			if votes[i] > votes[maxIdx] {
				maxIdx = i
			}
		}
		if maxIdx == 0 {
			break
		}
		votes[0]++
		votes[maxIdx]--
		count++
	}
	fmt.Println(count)
}
