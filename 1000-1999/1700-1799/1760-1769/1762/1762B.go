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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		v := make([]struct {
			val int64
			idx int
		}, n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			v[i].val = x
			v[i].idx = i
		}
		sort.Slice(v, func(i, j int) bool {
			return v[i].val < v[j].val
		})
		// Collect operations
		ops := make([][2]int64, 0, n)
		for i := 1; i < n; i++ {
			prev := v[i-1].val
			cur := v[i].val
			add := prev - (cur % prev)
			if add < 0 {
				add = 0
			}
			ops = append(ops, [2]int64{int64(v[i].idx + 1), add})
			v[i].val = cur + add
		}
		// Output
		fmt.Fprintln(writer, len(ops))
		for _, op := range ops {
			fmt.Fprintln(writer, op[0], op[1])
		}
	}
}
