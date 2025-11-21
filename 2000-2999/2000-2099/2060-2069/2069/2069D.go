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
		var s string
		fmt.Fscan(in, &s)
		n := len(s)

		pref := make([][26]int, n+1)
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i]
			pref[i+1][s[i]-'a']++
		}

		half := n / 2
		mismatch := make([]bool, half)
		prefBad := make([]int, n+1)
		suffDiff := make([]int, n+1)
		mismatchCount := 0

		for i := 0; i < half; i++ {
			j := n - 1 - i
			if s[i] != s[j] {
				mismatch[i] = true
				prefBad[j+1]++
				suffDiff[0]++
				suffDiff[i+1]--
				mismatchCount++
			}
		}

		if mismatchCount == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		for i := 1; i <= n; i++ {
			prefBad[i] += prefBad[i-1]
		}

		suffBad := make([]int, n+1)
		run := 0
		for i := 0; i <= n; i++ {
			run += suffDiff[i]
			suffBad[i] = run
		}

		maxJBefore := make([]int, n+1)
		maxJ := -1
		for L := 0; L <= n; L++ {
			if L > 0 {
				idx := L - 1
				if idx < half && mismatch[idx] {
					j := n - 1 - idx
					if j > maxJ {
						maxJ = j
					}
				}
			}
			maxJBefore[L] = maxJ
		}

		checkWindow := func(L, R int) bool {
			if prefBad[L] > 0 || suffBad[R] > 0 {
				return false
			}
			if maxJBefore[L] >= R {
				return false
			}

			rightLen := n - R
			startLeft := rightLen
			endLeft := L
			if n-L < endLeft {
				endLeft = n - L
			}

			startRight := R
			if rightLen > startRight {
				startRight = rightLen
			}
			endRight := n - L

			leftValid := startLeft < endLeft
			rightValid := startRight < endRight

			for c := 0; c < 26; c++ {
				sub := pref[R][c] - pref[L][c]
				dem := 0
				if leftValid {
					dem += pref[endLeft][c] - pref[startLeft][c]
				}
				if rightValid {
					dem += pref[endRight][c] - pref[startRight][c]
				}
				if dem > sub {
					return false
				}
				if (sub-dem)&1 == 1 {
					return false
				}
			}
			return true
		}

		exists := func(length int) bool {
			for L := 0; L+length <= n; L++ {
				if checkWindow(L, L+length) {
					return true
				}
			}
			return false
		}

		lo, hi := 0, n
		ans := n
		for lo <= hi {
			mid := (lo + hi) / 2
			if exists(mid) {
				ans = mid
				hi = mid - 1
			} else {
				lo = mid + 1
			}
		}

		fmt.Fprintln(out, ans)
	}
}
