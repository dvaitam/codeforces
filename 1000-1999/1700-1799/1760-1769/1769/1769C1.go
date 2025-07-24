package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(a []int) int {
	c := make([]int, 102)
	for _, v := range a {
		c[v]++
	}
	dp := [2]int{0, -1 << 60}
	best := 0
	for day := 1; day <= 101; day++ {
		nd := [2]int{-1 << 60, -1 << 60}
		for carry := 0; carry <= 1; carry++ {
			cur := dp[carry]
			if cur < 0 {
				continue
			}
			avail := c[day]
			if carry == 1 {
				newCarry := 0
				if avail > 0 {
					newCarry = 1
				}
				if cur+1 > nd[newCarry] {
					nd[newCarry] = cur + 1
				}
				if cur+1 > best {
					best = cur + 1
				}
				if avail > 0 {
					newCarry = 0
					if avail-1 > 0 {
						newCarry = 1
					}
					if cur+1 > nd[newCarry] {
						nd[newCarry] = cur + 1
					}
					if cur+1 > best {
						best = cur + 1
					}
				}
				newCarry = 0
				if avail > 0 {
					newCarry = 1
				}
				if 0 > nd[newCarry] {
					nd[newCarry] = 0
				}
			} else {
				if avail > 0 {
					newCarry := 0
					if avail-1 > 0 {
						newCarry = 1
					}
					if cur+1 > nd[newCarry] {
						nd[newCarry] = cur + 1
					}
					if cur+1 > best {
						best = cur + 1
					}
				}
				newCarry := 0
				if avail > 0 {
					newCarry = 1
				}
				if 0 > nd[newCarry] {
					nd[newCarry] = 0
				}
			}
		}
		dp = nd
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		res := solve(a)
		fmt.Fprintln(writer, res)
	}
}
