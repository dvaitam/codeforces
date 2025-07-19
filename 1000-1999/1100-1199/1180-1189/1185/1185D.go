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
	a := make([]pair, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i].val)
		a[i].idx = i + 1
	}
	if n <= 3 {
		fmt.Println(1)
		return
	}
	sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })
	// Scenario 1: remove one element in middle
	pos := -1
	ok := true
	if (a[n-1].val-a[0].val)%(n-2) == 0 {
		d := (a[n-1].val - a[0].val) / (n - 2)
		cnt := 0
		first := true
		for i := 0; i < n-1; i++ {
			expected := a[0].val + cnt*d
			if a[i].val != expected {
				if first {
					first = false
					pos = a[i].idx
					continue
				}
				ok = false
				break
			}
			cnt++
		}
	} else {
		ok = false
	}
	if ok && pos != -1 {
		fmt.Println(pos)
		return
	}
	// Scenario 2: remove last element
	ok = true
	if (a[n-2].val-a[0].val)%(n-2) == 0 {
		d := (a[n-2].val - a[0].val) / (n - 2)
		for i := 0; i < n-1; i++ {
			if a[i].val != a[0].val+i*d {
				ok = false
				break
			}
		}
	} else {
		ok = false
	}
	if ok {
		// removed element is the last one
		fmt.Println(a[n-1].idx)
		return
	}
	// Scenario 3: remove first element
	ok = true
	if (a[n-1].val-a[1].val)%(n-2) == 0 {
		d := (a[n-1].val - a[1].val) / (n - 2)
		for i := 0; i < n-2; i++ {
			if a[i+1].val != a[1].val+i*d {
				ok = false
				break
			}
		}
	} else {
		ok = false
	}
	if ok {
		fmt.Println(a[0].idx)
	} else {
		fmt.Println(-1)
	}
}

type pair struct {
	val int
	idx int
}
