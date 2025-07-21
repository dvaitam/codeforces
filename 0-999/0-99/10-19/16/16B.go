package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	type pair struct { b, a int64 }
	arr := make([]pair, m)
	for i := int64(0); i < m; i++ {
		var ai, bi int64
		fmt.Fscan(in, &ai, &bi)
		arr[i] = pair{b: bi, a: ai}
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].b > arr[j].b
	})
	var total int64
	remaining := n
	for _, p := range arr {
		if remaining == 0 {
			break
		}
		take := p.a
		if take > remaining {
			take = remaining
		}
		total += take * p.b
		remaining -= take
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprint(out, total)
}
