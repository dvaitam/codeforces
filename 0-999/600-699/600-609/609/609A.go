package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &m); err != nil {
		return
	}
	drives := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &drives[i])
	}
	sort.Slice(drives, func(i, j int) bool { return drives[i] > drives[j] })
	sum := 0
	cnt := 0
	for _, v := range drives {
		sum += v
		cnt++
		if sum >= m {
			break
		}
	}
	fmt.Fprintln(out, cnt)
}
