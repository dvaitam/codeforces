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
		fmt.Fscan(reader, &n)
		m := 2 * n
		d := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &d[i])
		}
		sort.Slice(d, func(i, j int) bool { return d[i] > d[j] })
		ok := true
		sum := int64(0)
		seen := make(map[int64]bool)
		for i := 0; i < m; i += 2 {
			if d[i] != d[i+1] || d[i]%2 == 1 {
				ok = false
				break
			}
			rem := d[i] - 2*sum
			k := int64(n - i/2)
			if rem <= 0 || rem%(2*k) != 0 {
				ok = false
				break
			}
			x := rem / (2 * k)
			if seen[x] || x <= 0 {
				ok = false
				break
			}
			seen[x] = true
			sum += x
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
