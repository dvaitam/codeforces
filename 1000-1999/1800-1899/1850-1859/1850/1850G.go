package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		mX := make(map[int]int)
		mY := make(map[int]int)
		mD1 := make(map[int]int)
		mD2 := make(map[int]int)
		for i := 0; i < n; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			mX[x]++
			mY[y]++
			mD1[x-y]++
			mD2[x+y]++
		}
		var ans int64
		for _, v := range mX {
			ans += int64(v * (v - 1))
		}
		for _, v := range mY {
			ans += int64(v * (v - 1))
		}
		for _, v := range mD1 {
			ans += int64(v * (v - 1))
		}
		for _, v := range mD2 {
			ans += int64(v * (v - 1))
		}
		fmt.Fprintln(out, ans)
	}
}
