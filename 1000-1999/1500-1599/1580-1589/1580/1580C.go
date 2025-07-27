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

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	x := make([]int, n+1)
	y := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &x[i], &y[i])
	}

	const B = 700
	// small[p][r] = number of trains with period p that are in maintenance when day%p==r
	small := make([][]int, B+1)
	for i := range small {
		small[i] = make([]int, B+1)
	}
	diff := make([]int, m+2) // difference array for large periods
	addDay := make([]int, n+1)
	active := make([]bool, n+1)
	curLarge := 0

	for day := 1; day <= m; day++ {
		curLarge += diff[day]
		var op, k int
		fmt.Fscan(in, &op, &k)
		if op == 1 {
			active[k] = true
			addDay[k] = day
			p := x[k] + y[k]
			if p <= B {
				off := (day + x[k]) % p
				for t := 0; t < y[k]; t++ {
					small[p][(off+t)%p]++
				}
			} else {
				for start := day + x[k]; start <= m; start += p {
					l := start
					r := start + y[k]
					if r > m+1 {
						r = m + 1
					}
					diff[l]++
					diff[r]--
				}
			}
		} else {
			p := x[k] + y[k]
			d := addDay[k]
			if p <= B {
				off := (d + x[k]) % p
				for t := 0; t < y[k]; t++ {
					small[p][(off+t)%p]--
				}
			} else {
				for start := d + x[k]; start <= m; start += p {
					l := start
					r := start + y[k]
					if r > m+1 {
						r = m + 1
					}
					if r <= day {
						continue
					}
					if l >= day {
						diff[l]--
						diff[r]++
					} else {
						diff[day]--
						diff[r]++
						curLarge--
					}
				}
			}
			active[k] = false
		}

		smallCnt := 0
		for p := 1; p <= B; p++ {
			smallCnt += small[p][day%p]
		}
		fmt.Fprintln(out, smallCnt+curLarge)
	}
}
