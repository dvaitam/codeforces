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
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		total += a[i]
	}
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
	if total <= b[0]+b[1] {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
