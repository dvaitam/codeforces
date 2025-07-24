package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func LISLength(arr []int) int {
	d := make([]int, 0, len(arr))
	for _, v := range arr {
		i := sort.Search(len(d), func(i int) bool { return d[i] >= v })
		if i == len(d) {
			d = append(d, v)
		} else {
			d[i] = v
		}
	}
	return len(d)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		pos := make([]int, n+1)
		for i, v := range a {
			pos[v] = i
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
			b[i] = pos[b[i]]
		}
		lis := LISLength(b)
		fmt.Fprintln(writer, n-lis)
	}
}
