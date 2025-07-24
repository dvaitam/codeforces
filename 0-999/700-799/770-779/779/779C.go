package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	idx := make([]int, n)
	for i := range idx {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		di := a[idx[i]] - b[idx[i]]
		dj := a[idx[j]] - b[idx[j]]
		return di < dj
	})
	var total int64
	for i, id := range idx {
		if i < k || a[id] <= b[id] {
			total += int64(a[id])
		} else {
			total += int64(b[id])
		}
	}
	fmt.Println(total)
}
