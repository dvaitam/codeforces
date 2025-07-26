package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func sumSlice(arr []int) int64 {
	var s int64
	for _, v := range arr {
		s += int64(v)
	}
	return s
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		var k int64
		fmt.Fscan(reader, &n, &m, &k)
		a := make([]int, n)
		b := make([]int, m)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &b[i])
		}
		sort.Ints(a)
		sort.Ints(b)

		// step 1: Jellyfish's move
		if b[m-1] > a[0] {
			a[0], b[m-1] = b[m-1], a[0]
			sort.Ints(a)
			sort.Ints(b)
		}
		sumAfterFirst := sumSlice(a)

		// step 2: Gellyfish's move
		if b[0] < a[n-1] {
			a[n-1], b[0] = b[0], a[n-1]
			sort.Ints(a)
			sort.Ints(b)
		}
		sumAfterSecond := sumSlice(a)

		if k%2 == 1 {
			fmt.Fprintln(writer, sumAfterFirst)
		} else {
			fmt.Fprintln(writer, sumAfterSecond)
		}
	}
}
