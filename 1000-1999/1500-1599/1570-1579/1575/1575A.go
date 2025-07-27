package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	titles := make([]struct {
		key string
		idx int
	}, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		b := []byte(s)
		for j := 1; j < m; j += 2 {
			b[j] = 'Z' - (b[j] - 'A')
		}
		titles[i].key = string(b)
		titles[i].idx = i + 1
	}
	sort.Slice(titles, func(i, j int) bool {
		return titles[i].key < titles[j].key
	})
	out := bufio.NewWriter(os.Stdout)
	for i, t := range titles {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, t.idx)
	}
	fmt.Fprintln(out)
	out.Flush()
}
