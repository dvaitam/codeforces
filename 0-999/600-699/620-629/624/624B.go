package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })

	var last int64 = 1 << 60
	var res int64
	for _, v := range a {
		if last == 0 {
			break
		}
		if v >= last {
			v = last - 1
		}
		if v < 0 {
			v = 0
		}
		res += v
		last = v
	}
	fmt.Println(res)
}
