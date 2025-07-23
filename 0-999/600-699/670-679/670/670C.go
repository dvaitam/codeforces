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
	freq := make(map[int]int, n)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		freq[x]++
	}
	var m int
	fmt.Fscan(in, &m)
	audio := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &audio[i])
	}
	sub := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &sub[i])
	}

	bestIdx := 1
	bestAudio := freq[audio[0]]
	bestSub := freq[sub[0]]
	for i := 1; i < m; i++ {
		a := freq[audio[i]]
		s := freq[sub[i]]
		if a > bestAudio || (a == bestAudio && s > bestSub) {
			bestAudio = a
			bestSub = s
			bestIdx = i + 1
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, bestIdx)
	out.Flush()
}
