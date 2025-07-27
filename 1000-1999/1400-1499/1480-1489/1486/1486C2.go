package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func query(l, r int) int {
	fmt.Fprintf(writer, "? %d %d\n", l, r)
	writer.Flush()
	var idx int
	fmt.Fscan(reader, &idx)
	return idx
}

func main() {
	defer writer.Flush()
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	s := query(1, n)

	// Determine which side contains the maximum element
	left := 1
	right := n

	if s > 1 {
		x := query(1, s)
		if x == s {
			right = s - 1
		} else {
			left = s + 1
		}
	} else {
		left = 2
	}

	if left <= right {
		if right < s { // search on the left side
			for left < right {
				mid := (left + right + 1) / 2
				if query(mid, s) == s {
					left = mid
				} else {
					right = mid - 1
				}
			}
		} else { // search on the right side
			for left < right {
				mid := (left + right) / 2
				if query(s, mid) == s {
					right = mid
				} else {
					left = mid + 1
				}
			}
		}
	}

	fmt.Fprintf(writer, "! %d\n", left)
	writer.Flush()
}
