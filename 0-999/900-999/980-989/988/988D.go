package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	m := make(map[int64]bool, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		m[a[i]] = true
	}
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	// remove duplicates in sorted slice
	uniq := a[:0]
	for _, v := range a {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	a = uniq
	last := a[len(a)-1]
	found2 := false
	var ansX, ansY int64
	for _, x := range a {
		var prevDiff int64 = -1
		for diff := int64(1); x+diff <= last; diff <<= 1 {
			if m[x+diff] {
				if prevDiff != -1 && diff == prevDiff*2 {
					// found triple
					fmt.Println(3)
					fmt.Printf("%d %d %d\n", x, x+prevDiff, x+diff)
					return
				}
				prevDiff = diff
			}
		}
		if prevDiff != -1 && !found2 {
			found2 = true
			ansX = x
			ansY = x + prevDiff
		}
	}
	if found2 {
		fmt.Println(2)
		fmt.Printf("%d %d\n", ansX, ansY)
	} else {
		fmt.Println(1)
		fmt.Printf("%d\n", a[0])
	}
}
