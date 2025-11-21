package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	n := len(s)
	prefix := make([][3]int, n+1)
	for i, ch := range s {
		prefix[i+1] = prefix[i]
		switch ch {
		case 'G':
			prefix[i+1][0]++
		case 'S':
			prefix[i+1][1]++
		case 'B':
			prefix[i+1][2]++
		}
	}

	var q int
	fmt.Fscan(in, &q)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(in, &l, &r)
		counts := [3]int{
			prefix[r][0] - prefix[l-1][0],
			prefix[r][1] - prefix[l-1][1],
			prefix[r][2] - prefix[l-1][2],
		}

		// Sort counts ascending using three comparisons.
		if counts[0] > counts[1] {
			counts[0], counts[1] = counts[1], counts[0]
		}
		if counts[1] > counts[2] {
			counts[1], counts[2] = counts[2], counts[1]
		}
		if counts[0] > counts[1] {
			counts[0], counts[1] = counts[1], counts[0]
		}

		fmt.Fprintln(out, counts[2]+counts[0])
	}
}
