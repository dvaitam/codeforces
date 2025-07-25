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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)

		divs := divisors(n)
		ans := n
		for _, d := range divs {
			if check(s, d) {
				ans = d
				break
			}
		}
		fmt.Fprintln(out, ans)
	}
}

func divisors(n int) []int {
	ds := make([]int, 0)
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			ds = append(ds, i)
			if i*i != n {
				ds = append(ds, n/i)
			}
		}
	}
	// sort the divisors
	for i := 0; i < len(ds); i++ {
		for j := i + 1; j < len(ds); j++ {
			if ds[j] < ds[i] {
				ds[i], ds[j] = ds[j], ds[i]
			}
		}
	}
	return ds
}

func check(s string, d int) bool {
	n := len(s)
	freq := make([][26]int, d)
	for i := 0; i < n; i++ {
		freq[i%d][s[i]-'a']++
	}
	block := n / d
	mism := 0
	for i := 0; i < d; i++ {
		maxC := 0
		for j := 0; j < 26; j++ {
			if freq[i][j] > maxC {
				maxC = freq[i][j]
			}
		}
		mism += block - maxC
		if mism > 1 {
			return false
		}
	}
	return mism <= 1
}
