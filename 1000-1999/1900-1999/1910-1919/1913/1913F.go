package main

import (
	"bufio"
	"fmt"
	"os"
)

func countPal(s []byte) int64 {
	n := len(s)
	rad1 := make([]int, n)
	l, r := 0, -1
	for i := 0; i < n; i++ {
		k := 1
		if i <= r {
			if rad1[l+r-i] < r-i+1 {
				k = rad1[l+r-i]
			} else {
				k = r - i + 1
			}
		}
		for i-k >= 0 && i+k < n && s[i-k] == s[i+k] {
			k++
		}
		rad1[i] = k
		if i+k-1 > r {
			l = i - k + 1
			r = i + k - 1
		}
	}
	rad2 := make([]int, n)
	l, r = 0, -1
	for i := 0; i < n; i++ {
		k := 0
		if i <= r {
			if rad2[l+r-i+1] < r-i+1 {
				k = rad2[l+r-i+1]
			} else {
				k = r - i + 1
			}
		}
		for i-k-1 >= 0 && i+k < n && s[i-k-1] == s[i+k] {
			k++
		}
		rad2[i] = k
		if i+k-1 > r {
			l = i - k
			r = i + k - 1
		}
	}
	var count int64
	for i := 0; i < n; i++ {
		count += int64(rad1[i])
		count += int64(rad2[i])
	}
	return count
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	bs := []byte(s)

	bestCount := countPal(bs)
	bestStr := s

	for i := 0; i < n; i++ {
		orig := bs[i]
		for c := byte('a'); c <= byte('z'); c++ {
			if c == orig {
				continue
			}
			bs[i] = c
			cnt := countPal(bs)
			str := string(bs)
			if cnt > bestCount || (cnt == bestCount && str < bestStr) {
				bestCount = cnt
				bestStr = str
			}
		}
		bs[i] = orig
	}

	fmt.Fprintln(out, bestCount)
	fmt.Fprintln(out, bestStr)
}
