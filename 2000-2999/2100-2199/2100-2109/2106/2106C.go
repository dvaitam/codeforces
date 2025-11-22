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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		hasFixed := false
		var needX int64
		ok := true
		for i := 0; i < n; i++ {
			if b[i] != -1 {
				if !hasFixed {
					hasFixed = true
					needX = a[i] + b[i]
				} else if needX != a[i]+b[i] {
					ok = false
					break
				}
			}
		}

		var ans int64
		if ok {
			if hasFixed {
				x := needX
				for i := 0; i < n; i++ {
					val := x - a[i]
					if val < 0 || val > k {
						ok = false
						break
					}
					if b[i] != -1 && b[i] != val {
						ok = false
						break
					}
				}
				if ok {
					ans = 1
				} else {
					ans = 0
				}
			} else {
				// intersection of intervals [a[i], a[i]+k]
				L := a[0]
				R := a[0] + k
				for i := 1; i < n; i++ {
					if a[i] > L {
						L = a[i]
					}
					if a[i]+k < R {
						R = a[i] + k
					}
				}
				if L > R {
					ans = 0
				} else {
					ans = R - L + 1
				}
			}
		} else {
			ans = 0
		}

		fmt.Fprintln(out, ans)
	}
}
