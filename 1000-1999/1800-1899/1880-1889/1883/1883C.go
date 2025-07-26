// Solution for Codeforces problem 1883C
// Minimum increments to make product divisible by k (k<=5).
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}
		ans := 0
		switch k {
		case 2:
			ans = 1
			for _, v := range a {
				if v%2 == 0 {
					ans = 0
					break
				}
			}
		case 3:
			ans = 2
			for _, v := range a {
				d := (3 - v%3) % 3
				if d < ans {
					ans = d
					if ans == 0 {
						break
					}
				}
			}
		case 4:
			twos := 0
			for _, v := range a {
				x := v
				for x%2 == 0 {
					twos++
					x /= 2
					if twos >= 2 {
						break
					}
				}
				if twos >= 2 {
					break
				}
			}
			if twos >= 2 {
				ans = 0
				break
			}
			best4 := int(1e9)
			min1, min2 := int(1e9), int(1e9)
			for _, v := range a {
				c4 := (4 - v%4) % 4
				if c4 < best4 {
					best4 = c4
				}
				ce := 0
				if v%2 != 0 {
					ce = 1
				}
				if ce < min1 {
					min2 = min1
					min1 = ce
				} else if ce < min2 {
					min2 = ce
				}
			}
			if min2 == int(1e9) {
				min2 = 0
			}
			if best4 < min1+min2 {
				ans = best4
			} else {
				ans = min1 + min2
			}
		case 5:
			ans = 4
			for _, v := range a {
				d := (5 - v%5) % 5
				if d < ans {
					ans = d
					if ans == 0 {
						break
					}
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
