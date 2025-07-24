package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var h int64
	fmt.Fscan(in, &n, &h)
	x := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i])
	}
	var last int64 = -1 << 60
	var cnt int64
	for _, v := range x {
		if last >= v-1 {
			continue
		}
		last = v + 1
		cnt++
	}
	ans := cnt * h
	fmt.Println(ans)
}
