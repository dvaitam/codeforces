package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	vals := make([]int, 5)
	freq := make(map[int]int)
	sum := 0
	for i := 0; i < 5; i++ {
		fmt.Fscan(in, &vals[i])
		sum += vals[i]
		freq[vals[i]]++
	}
	best := 0
	for v, c := range freq {
		if c >= 2 {
			cand := v * 2
			if c >= 3 {
				cand = v * 3
			}
			if cand > best {
				best = cand
			}
		}
	}
	fmt.Println(sum - best)
}
