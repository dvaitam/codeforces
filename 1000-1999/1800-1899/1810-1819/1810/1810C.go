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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var c, d int64
		fmt.Fscan(reader, &n, &c, &d)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		sort.Ints(arr)
		uniq := make([]int, 0, n)
		last := -1
		for _, v := range arr {
			if v != last {
				uniq = append(uniq, v)
				last = v
			}
		}
		k := len(uniq)

		cost := func(m int) int64 {
			p := sort.Search(len(uniq), func(i int) bool { return uniq[i] > m })
			del := n - p
			ins := m - p
			return int64(del)*c + int64(ins)*d
		}

		best := cost(1)
		for i := 0; i < k; i++ {
			v := uniq[i]
			cur := cost(v)
			if cur < best {
				best = cur
			}
		}
		fmt.Fprintln(writer, best)
	}
}
