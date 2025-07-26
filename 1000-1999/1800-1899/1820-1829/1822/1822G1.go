package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxA = 1000000

func divisors(x int) []int {
	res := []int{}
	for d := 1; d*d <= x; d++ {
		if x%d == 0 {
			res = append(res, d)
			if d != x/d {
				res = append(res, x/d)
			}
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	freq := make([]int, maxA+1)
	used := make([]int, 0)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		used = used[:0]
		for i := 0; i < n; i++ {
			var v int
			fmt.Fscan(reader, &v)
			if freq[v] == 0 {
				used = append(used, v)
			}
			freq[v]++
		}

		var ans int64
		for _, x := range used {
			cx := freq[x]
			divs := divisors(x)
			for _, b := range divs {
				if x*b > maxA {
					continue
				}
				y := x / b
				z := x * b
				cy := freq[y]
				cz := freq[z]
				if cy == 0 || cz == 0 {
					continue
				}
				if y == x && z == x {
					if cx >= 3 {
						ans += int64(cx) * int64(cx-1) * int64(cx-2)
					}
				} else if y == x {
					if cx >= 2 {
						ans += int64(cx) * int64(cx-1) * int64(cz)
					}
				} else if z == x {
					if cx >= 2 {
						ans += int64(cy) * int64(cx) * int64(cx-1)
					}
				} else {
					ans += int64(cx) * int64(cy) * int64(cz)
				}
			}
		}
		fmt.Fprintln(writer, ans)

		for _, v := range used {
			freq[v] = 0
		}
	}
}
