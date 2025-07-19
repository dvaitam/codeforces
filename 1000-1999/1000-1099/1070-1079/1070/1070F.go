package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([][]int64, 4)
	for i := 0; i < n; i++ {
		var s string
		var x int64
		fmt.Fscan(reader, &s, &x)
		idx := int(s[0]-'0')*2 + int(s[1]-'0')
		a[idx] = append(a[idx], x)
	}
	m := min(len(a[1]), len(a[2]))
	for i := 1; i <= 2; i++ {
		sort.Slice(a[i], func(p, q int) bool { return a[i][p] < a[i][q] })
		xsize := len(a[i]) - m
		a[0] = append(a[0], a[i][:xsize]...)
		a[3] = append(a[3], a[i][xsize:]...)
	}
	sort.Slice(a[0], func(p, q int) bool { return a[0][p] < a[0][q] })
	k := min(len(a[0]), len(a[3])-2*m)
	if k > 0 {
		a[3] = append(a[3], a[0][len(a[0])-k:]...)
	}
	var sum int64
	for _, v := range a[3] {
		sum += v
	}
	fmt.Fprintln(writer, sum)
}
