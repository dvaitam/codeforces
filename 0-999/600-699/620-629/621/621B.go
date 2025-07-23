package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	const offset = 1000
	var diag1 [2001]int64
	var diag2 [2001]int64
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		diag1[x-y+offset]++
		diag2[x+y]++
	}
	var ans int64
	for _, c := range diag1 {
		ans += c * (c - 1) / 2
	}
	for _, c := range diag2 {
		ans += c * (c - 1) / 2
	}
	fmt.Println(ans)
}
