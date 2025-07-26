package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	maxA := 0
	for i := range a {
		fmt.Fscan(in, &a[i])
		if a[i] > maxA {
			maxA = a[i]
		}
	}

	diff := make([]int64, maxA+2)
	add := func(l, r int, v int64) {
		if l > r || l > maxA {
			return
		}
		if r > maxA {
			r = maxA
		}
		diff[l] += v
		diff[r+1] -= v
	}

	processSingle := func(x int) {
		k := 1
		for k <= maxA {
			q := (x + k - 1) / k
			var next int
			if q == 1 {
				next = maxA
			} else {
				next = (x - 1) / (q - 1)
				if next > maxA {
					next = maxA
				}
			}
			add(k, next, int64(q))
			if q == 1 {
				break
			}
			k = next + 1
		}
	}

	processPair := func(x, y int) {
		k := 1
		for k <= maxA {
			qx := (x + k - 1) / k
			qy := (y + k - 1) / k
			val := qy - qx
			if val < 0 {
				val = 0
			}
			var nextX, nextY int
			if qx == 1 {
				nextX = maxA
			} else {
				nextX = (x - 1) / (qx - 1)
				if nextX > maxA {
					nextX = maxA
				}
			}
			if qy == 1 {
				nextY = maxA
			} else {
				nextY = (y - 1) / (qy - 1)
				if nextY > maxA {
					nextY = maxA
				}
			}
			end := nextX
			if nextY < end {
				end = nextY
			}
			if end > maxA {
				end = maxA
			}
			add(k, end, int64(val))
			if qx == 1 && qy == 1 {
				break
			}
			k = end + 1
		}
	}

	processSingle(a[0])
	for i := 1; i < n; i++ {
		processPair(a[i-1], a[i])
	}

	ans := make([]int64, maxA+1)
	cur := int64(0)
	for i := 1; i <= maxA; i++ {
		cur += diff[i]
		ans[i] = cur
	}

	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= maxA; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
	out.Flush()
}
