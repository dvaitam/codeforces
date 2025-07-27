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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		x := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &x[i])
		}

		pref := make([]int64, n)
		var maxPref int64
		for i := 0; i < n; i++ {
			if i == 0 {
				pref[i] = a[i]
			} else {
				pref[i] = pref[i-1] + a[i]
			}
			if pref[i] > maxPref {
				maxPref = pref[i]
			}
		}
		total := pref[n-1]

		prefixMax := make([]int64, n)
		prefixMax[0] = pref[0]
		for i := 1; i < n; i++ {
			if pref[i] > prefixMax[i-1] {
				prefixMax[i] = pref[i]
			} else {
				prefixMax[i] = prefixMax[i-1]
			}
		}

		for i := 0; i < m; i++ {
			q := x[i]
			if maxPref >= q {
				l, r := 0, n-1
				for l < r {
					mid := (l + r) / 2
					if prefixMax[mid] >= q {
						r = mid
					} else {
						l = mid + 1
					}
				}
				if i+1 == m {
					fmt.Fprintln(writer, l)
				} else {
					fmt.Fprint(writer, l, " ")
				}
				continue
			}

			if total <= 0 {
				if i+1 == m {
					fmt.Fprintln(writer, -1)
				} else {
					fmt.Fprint(writer, -1, " ")
				}
				continue
			}

			loops := (q - maxPref + total - 1) / total
			remain := q - loops*total
			// remain <= maxPref, and remain may be <=0
			if remain <= 0 {
				// reached within loops cycles at or before returning to start
				res := loops * int64(n)
				if i+1 == m {
					fmt.Fprintln(writer, res)
				} else {
					fmt.Fprint(writer, res, " ")
				}
				continue
			}
			l, r := 0, n-1
			for l < r {
				mid := (l + r) / 2
				if prefixMax[mid] >= remain {
					r = mid
				} else {
					l = mid + 1
				}
			}
			res := loops*int64(n) + int64(l)
			if i+1 == m {
				fmt.Fprintln(writer, res)
			} else {
				fmt.Fprint(writer, res, " ")
			}
		}
	}
}
