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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		has0 := make([]bool, n+2)
		has1 := make([]bool, n+2)
		var xs0, xs1 []int

		for i := 0; i < n; i++ {
			var x, y int
			fmt.Fscan(in, &x, &y)
			if y == 0 {
				if !has0[x] {
					xs0 = append(xs0, x)
				}
				has0[x] = true
			} else {
				if !has1[x] {
					xs1 = append(xs1, x)
				}
				has1[x] = true
			}
		}

		common := 0
		for x := 0; x <= n; x++ {
			if has0[x] && has1[x] {
				common++
			}
		}

		count0 := len(xs0)
		count1 := len(xs1)
		part0 := 0
		if count0 > 0 {
			part0 = count0 - 1
		}
		part1 := 0
		if count1 > 0 {
			part1 = count1 - 1
		}

		ans := int64(common) * int64(part0+part1)

		for _, x := range xs1 {
			if x > 0 && x < n {
				if has0[x-1] && has0[x+1] {
					ans++
				}
			}
		}
		for _, x := range xs0 {
			if x > 0 && x < n {
				if has1[x-1] && has1[x+1] {
					ans++
				}
			}
		}

		fmt.Fprintln(out, ans)
	}
}
