package main

import (
	"bufio"
	"fmt"
	"os"
)

func minCrystals(h int, p []int) int {
	p = append(p, 0) // ground level
	idx := 1         // index of next platform below current height
	cur := h
	ans := 0
	for cur > 2 {
		// skip all platforms not below current height
		for idx < len(p) && p[idx] >= cur {
			idx++
		}
		if idx >= len(p) {
			break
		}
		if p[idx] == cur-1 {
			// platform right below will disappear after toggle
			next := 0
			if idx+1 < len(p) {
				next = p[idx+1]
			}
			if cur-next > 2 {
				// need a crystal to create platform at cur-2
				ans++
				cur -= 2
				for idx < len(p) && p[idx] >= cur {
					idx++
				}
			} else {
				// safely fall to the next existing platform
				cur = next
				idx += 2
			}
		} else {
			// no platform at cur-1; it will appear after toggle
			cur--
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	for ; q > 0; q-- {
		var h, n int
		fmt.Fscan(in, &h, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		fmt.Fprintln(out, minCrystals(h, p))
	}
}
