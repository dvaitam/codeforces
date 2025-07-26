package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}
		sort.Ints(a)
		sort.Ints(b)
		l, r := 0, m-1
		i, j := 0, n-1
		var sum int64
		for i <= j {
			diff1 := abs(a[i] - b[r])
			diff2 := abs(a[j] - b[l])
			if diff1 >= diff2 {
				sum += int64(diff1)
				i++
				r--
			} else {
				sum += int64(diff2)
				j--
				l++
			}
		}
		fmt.Fprintln(writer, sum)
	}
}
