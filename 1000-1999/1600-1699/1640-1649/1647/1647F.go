package main

import (
	"bufio"
	"fmt"
	"os"
)

func canSplit(a []int, pos1, pos2 int) bool {
	n := len(a)
	last1, last2 := -1, -1
	peak1, peak2 := a[pos1], a[pos2]
	for i := 0; i < n; i++ {
		val := a[i]
		if i == pos1 {
			last1 = val
			continue
		}
		if i == pos2 {
			last2 = val
			continue
		}
		ok1 := false
		if i < pos1 {
			if val > last1 && val < peak1 {
				ok1 = true
			}
		} else if i > pos1 {
			if val < last1 {
				ok1 = true
			}
		}
		ok2 := false
		if i < pos2 {
			if val > last2 && val < peak2 {
				ok2 = true
			}
		} else if i > pos2 {
			if val < last2 {
				ok2 = true
			}
		}
		if ok1 && !ok2 {
			last1 = val
		} else if ok2 && !ok1 {
			last2 = val
		} else if ok1 && ok2 {
			if last1 <= last2 {
				last1 = val
			} else {
				last2 = val
			}
		} else {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	pos := make(map[int]int)
	maxVal := -1
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		pos[a[i]] = i
		if a[i] > maxVal {
			maxVal = a[i]
		}
	}

	pairs := make(map[[2]int]struct{})
	posMax := pos[maxVal]
	for val, p := range pos {
		if val == maxVal {
			continue
		}
		if canSplit(a, posMax, p) {
			pair := [2]int{val, maxVal}
			if pair[0] > pair[1] {
				pair[0], pair[1] = pair[1], pair[0]
			}
			pairs[pair] = struct{}{}
		}
	}
	fmt.Fprintln(out, len(pairs))
}
