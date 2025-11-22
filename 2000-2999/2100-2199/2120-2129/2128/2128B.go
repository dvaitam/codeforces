package main

import (
	"bufio"
	"fmt"
	"os"
)

// runDir: 0 unknown/length 1, 1 increasing, -1 decreasing
// runLen: length of current monotone run ending at last chosen element.

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
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}

		l, r := 0, n-1
		last := 0
		runLen := 0
		runDir := 0
		ans := make([]byte, 0, n)

		for l <= r {
			choice := byte(0)
			var nLast, nRunLen, nRunDir int

			try := func(side byte, val int, nl, nr int) bool {
				dir := runDir
				lenRun := runLen
				if lenRun == 0 {
					lenRun = 1
					dir = 0
				} else if lenRun == 1 {
					lenRun = 2
					if val > last {
						dir = 1
					} else {
						dir = -1
					}
				} else {
					cmp := 1
					if val < last {
						cmp = -1
					}
					if cmp == dir {
						lenRun++
					} else {
						lenRun = 2
						dir = cmp
					}
				}

				if lenRun >= 5 {
					return false
				}

				if lenRun == 4 && nl <= nr {
					needLess := dir == 1 // run is increasing; need a smaller next value to break
					ok := false
					if needLess {
						if p[nl] < val || p[nr] < val {
							ok = true
						}
					} else { // decreasing run, need bigger next value
						if p[nl] > val || p[nr] > val {
							ok = true
						}
					}
					if !ok {
						return false
					}
				}

				choice = side
				nLast = val
				nRunLen = lenRun
				nRunDir = dir
				return true
			}

			// Prefer left if both work, but any valid choice is fine.
			if l == r {
				_ = try('L', p[l], l+1, r)
			} else {
				if !try('L', p[l], l+1, r) {
					try('R', p[r], l, r-1)
				}
			}

			// problem guarantees existence; choice must be set
			ans = append(ans, choice)
			last = nLast
			runLen = nRunLen
			runDir = nRunDir
			if choice == 'L' {
				l++
			} else {
				r--
			}
		}

		fmt.Fprintln(out, string(ans))
	}
}

