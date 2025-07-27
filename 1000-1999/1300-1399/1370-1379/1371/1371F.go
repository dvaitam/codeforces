package main

import (
	"bufio"
	"fmt"
	"os"
)

// computeHoles calculates the final hole index for each conveyor based on
// the current configuration string s. Hole indices range from 0 to len(s).
func computeHoles(s []byte) []int {
	n := len(s)
	hole := make([]int, n)
	left := make([]int, n)
	right := make([]int, n)

	// process '<' segments from left to right
	for i := 0; i < n; i++ {
		if s[i] == '<' {
			if i == 0 {
				left[i] = 0
			} else if s[i-1] == '<' {
				left[i] = left[i-1]
			} else {
				left[i] = i
			}
		}
	}

	// process '>' segments from right to left
	for i := n - 1; i >= 0; i-- {
		if s[i] == '>' {
			if i == n-1 {
				right[i] = n
			} else if s[i+1] == '>' {
				right[i] = right[i+1]
			} else {
				right[i] = i + 1
			}
		}
	}

	for i := 0; i < n; i++ {
		if s[i] == '<' {
			hole[i] = left[i]
		} else {
			hole[i] = right[i]
		}
	}
	return hole
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	var str string
	fmt.Fscan(reader, &str)
	s := []byte(str)

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		r--
		// flip the conveyors in the range [l,r]
		for i := l; i <= r; i++ {
			if s[i] == '<' {
				s[i] = '>'
			} else {
				s[i] = '<'
			}
		}

		hole := computeHoles(s)
		count := make(map[int]int)
		maxCnt := 0
		for i := l; i <= r; i++ {
			h := hole[i]
			count[h]++
			if count[h] > maxCnt {
				maxCnt = count[h]
			}
		}
		fmt.Fprintln(writer, maxCnt)
	}
}
