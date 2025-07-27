package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	both := make([]int, 0)
	aOnly := make([]int, 0)
	bOnly := make([]int, 0)
	for i := 0; i < n; i++ {
		var t, a, b int
		fmt.Fscan(reader, &t, &a, &b)
		if a == 1 && b == 1 {
			both = append(both, t)
		} else if a == 1 {
			aOnly = append(aOnly, t)
		} else if b == 1 {
			bOnly = append(bOnly, t)
		}
	}
	sort.Ints(both)
	sort.Ints(aOnly)
	sort.Ints(bOnly)
	pb := make([]int, len(both)+1)
	for i := 1; i <= len(both); i++ {
		pb[i] = pb[i-1] + both[i-1]
	}
	pa := make([]int, len(aOnly)+1)
	for i := 1; i <= len(aOnly); i++ {
		pa[i] = pa[i-1] + aOnly[i-1]
	}
	pc := make([]int, len(bOnly)+1)
	for i := 1; i <= len(bOnly); i++ {
		pc[i] = pc[i-1] + bOnly[i-1]
	}
	const inf int64 = 1<<63 - 1
	ans := inf
	limit := k
	if len(both) < limit {
		limit = len(both)
	}
	for i := 0; i <= limit; i++ {
		need := k - i
		if need <= len(aOnly) && need <= len(bOnly) {
			t := int64(pb[i]) + int64(pa[need]) + int64(pc[need])
			if t < ans {
				ans = t
			}
		}
	}
	if ans == inf {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, ans)
	}
}
