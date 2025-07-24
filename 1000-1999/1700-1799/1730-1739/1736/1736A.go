package main

import (
	"bufio"
	"fmt"
	"os"
)

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
		b := make([]int, n)
		cntA, cntB := 0, 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			if a[i] == 1 {
				cntA++
			}
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
			if b[i] == 1 {
				cntB++
			}
		}
		mism := 0
		for i := 0; i < n; i++ {
			if a[i] != b[i] {
				mism++
			}
		}
		diff := cntA - cntB
		if diff < 0 {
			diff = -diff
		}
		ans := mism
		if diff+1 < ans {
			ans = diff + 1
		}
		fmt.Fprintln(writer, ans)
	}
}
