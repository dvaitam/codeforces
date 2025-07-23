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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	const days = 366
	male := make([]int, days+1)
	female := make([]int, days+1)

	for i := 0; i < n; i++ {
		var g string
		var a, b int
		fmt.Fscan(in, &g, &a, &b)
		if g == "M" {
			for d := a; d <= b; d++ {
				male[d]++
			}
		} else {
			for d := a; d <= b; d++ {
				female[d]++
			}
		}
	}

	ans := 0
	for d := 1; d <= days; d++ {
		if male[d] < female[d] {
			if ans < male[d]*2 {
				ans = male[d] * 2
			}
		} else {
			if ans < female[d]*2 {
				ans = female[d] * 2
			}
		}
	}
	fmt.Fprintln(out, ans)
}
