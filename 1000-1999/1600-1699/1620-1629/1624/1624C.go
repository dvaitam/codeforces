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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })
		used := make([]bool, n+1)
		ok := true
		for _, v := range a {
			for v > n || (v > 0 && used[v]) {
				v /= 2
			}
			if v == 0 {
				ok = false
				break
			}
			used[v] = true
		}
		if ok {
			for i := 1; i <= n; i++ {
				if !used[i] {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
