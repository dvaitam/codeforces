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

	var n, w int
	if _, err := fmt.Fscan(in, &n, &w); err != nil {
		return
	}

	ans := make([]int64, w+1)
	diff := make([]int64, w+2)

	for ; n > 0; n-- {
		var l int
		fmt.Fscan(in, &l)
		arr := make([]int64, l)
		for i := 0; i < l; i++ {
			fmt.Fscan(in, &arr[i])
		}
		if l == 0 {
			// nothing to add
			continue
		}
		// prefix maxima
		pref := make([]int64, l)
		pref[0] = arr[0]
		for i := 1; i < l; i++ {
			if arr[i] > pref[i-1] {
				pref[i] = arr[i]
			} else {
				pref[i] = pref[i-1]
			}
		}
		// suffix maxima
		suf := make([]int64, l)
		suf[l-1] = arr[l-1]
		for i := l - 2; i >= 0; i-- {
			if arr[i] > suf[i+1] {
				suf[i] = arr[i]
			} else {
				suf[i] = suf[i+1]
			}
		}
		maxVal := pref[l-1]
		if maxVal < 0 {
			maxVal = 0
		}

		prefixLimit := l
		if prefixLimit > w {
			prefixLimit = w
		}
		for j := 0; j < prefixLimit; j++ {
			v := pref[j]
			if v < 0 {
				v = 0
			}
			ans[j] += v
		}

		origSuffixStart := w - l
		suffixStart := origSuffixStart
		if suffixStart < 0 {
			suffixStart = 0
		}
		if suffixStart < prefixLimit {
			suffixStart = prefixLimit
		}
		for j := suffixStart; j < w; j++ {
			idx := j - origSuffixStart
			if idx >= l {
				ans[j] += 0
			} else {
				v := suf[idx]
				if v < 0 {
					v = 0
				}
				ans[j] += v
			}
		}

		left := prefixLimit
		right := origSuffixStart - 1
		if left <= right {
			diff[left] += maxVal
			diff[right+1] -= maxVal
		}
	}

	var add int64
	for i := 0; i < w; i++ {
		add += diff[i]
		ans[i] += add
	}

	for i := 0; i < w; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
