package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func beauty(arr []int) int64 {
	n := len(arr)
	b := append([]int(nil), arr...)
	sort.Ints(b)
	l, r := -1, -1
	for i := 0; i < n; i++ {
		if arr[i] != b[i] {
			if l == -1 {
				l = i
			}
			r = i
		}
	}
	if l == -1 {
		return 0
	}
	return int64(r - l)
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
		var ans int64
		for l := 0; l < n; l++ {
			for r := l; r < n; r++ {
				ans += beauty(a[l : r+1])
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
