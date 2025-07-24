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

	var n int
	var m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	a := make([]int, n)
	var sum int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		sum += int64(a[i])
	}
	sort.Ints(a)
	i, j := 0, n-1
	wait := 0
	for i < j {
		if a[i]+a[j] <= m {
			i++
			j--
		} else {
			wait++
			j--
		}
	}
	ans := sum + int64(wait) + 1
	fmt.Fprintln(writer, ans)
}
