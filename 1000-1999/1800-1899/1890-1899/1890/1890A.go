package main

import (
	"bufio"
	"fmt"
	"os"
)

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			freq[x]++
		}
		if len(freq) == 1 {
			fmt.Fprintln(out, "Yes")
			continue
		}
		if len(freq) != 2 {
			fmt.Fprintln(out, "No")
			continue
		}
		counts := make([]int, 0, 2)
		for _, c := range freq {
			counts = append(counts, c)
		}
		if absInt(counts[0]-counts[1]) <= 1 {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}
