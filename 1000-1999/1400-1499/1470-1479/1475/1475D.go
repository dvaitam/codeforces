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
		var m int64
		fmt.Fscan(reader, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		var total int64
		ones := []int64{}
		twos := []int64{}
		for i := 0; i < n; i++ {
			total += a[i]
			if b[i] == 1 {
				ones = append(ones, a[i])
			} else {
				twos = append(twos, a[i])
			}
		}
		if total < m {
			fmt.Fprintln(writer, -1)
			continue
		}
		sort.Slice(ones, func(i, j int) bool { return ones[i] > ones[j] })
		sort.Slice(twos, func(i, j int) bool { return twos[i] > twos[j] })
		pref1 := make([]int64, len(ones)+1)
		for i := 0; i < len(ones); i++ {
			pref1[i+1] = pref1[i] + ones[i]
		}
		pref2 := make([]int64, len(twos)+1)
		for i := 0; i < len(twos); i++ {
			pref2[i+1] = pref2[i] + twos[i]
		}
		ans := int(1e9)
		i := len(ones)
		for j := 0; j <= len(twos); j++ {
			mem := pref2[j]
			for i > 0 && mem+pref1[i-1] >= m {
				i--
			}
			if mem+pref1[i] >= m {
				cost := 2*j + i
				if cost < ans {
					ans = cost
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
