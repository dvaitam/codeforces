package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	mod1 int64 = 1000000007
	mod2 int64 = 1000000009
	base int64 = 911382323
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}
		ans := solve(a, n, k)
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}

func solve(a []int, n, k int) []int {
	if k == 0 {
		return []int{}
	}

	pow1 := make([]int64, n+1)
	pow2 := make([]int64, n+1)
	pow1[0], pow2[0] = 1, 1
	for i := 1; i <= n; i++ {
		pow1[i] = pow1[i-1] * base % mod1
		pow2[i] = pow2[i-1] * base % mod2
	}

	pref1 := make([]int64, n+1)
	pref2 := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref1[i+1] = (pref1[i]*base + int64(a[i])) % mod1
		pref2[i+1] = (pref2[i]*base + int64(a[i])) % mod2
	}

	rev := make([]int, n)
	for i := 0; i < n; i++ {
		rev[i] = a[n-1-i]
	}
	prefRev1 := make([]int64, n+1)
	prefRev2 := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefRev1[i+1] = (prefRev1[i]*base + int64(rev[i])) % mod1
		prefRev2[i+1] = (prefRev2[i]*base + int64(rev[i])) % mod2
	}

	isPal := func(l, r int) bool {
		if l > r {
			return true
		}
		hf1 := getHash(pref1, pow1, l, r, mod1)
		hb1 := getHash(prefRev1, pow1, n-1-r, n-1-l, mod1)
		if hf1 != hb1 {
			return false
		}
		hf2 := getHash(pref2, pow2, l, r, mod2)
		hb2 := getHash(prefRev2, pow2, n-1-r, n-1-l, mod2)
		return hf2 == hb2
	}

	palSuffix := make([]bool, n)
	for start := 0; start < n; start++ {
		palSuffix[start] = isPal(start, n-1)
	}

	// Values placed right before a palindromic suffix would immediately close a longer palindrome.
	forbid := make([]bool, n+1)
	forbid[a[n-1]] = true
	for start := 1; start < n; start++ {
		if palSuffix[start] {
			forbid[a[start-1]] = true
		}
	}

	curr := make([]int, n)
	copy(curr, a)
	ans := make([]int, 0, k)

	first := -1
	for val := 1; val <= n; val++ {
		if !forbid[val] {
			first = val
			break
		}
	}
	if first == -1 {
		first = 1
	}
	ans = append(ans, first)
	curr = append(curr, first)

	// Once the ending is palindrome-free, it suffices to avoid the last two values going forward.
	for i := 1; i < k; i++ {
		last := curr[len(curr)-1]
		second := curr[len(curr)-2]
		val := 1
		for val <= n {
			if val != last && val != second {
				break
			}
			val++
		}
		if val > n {
			val = 1
			for val <= n {
				if val != last && val != second {
					break
				}
				val++
			}
		}
		ans = append(ans, val)
		curr = append(curr, val)
	}
	return ans
}

func getHash(pref, pow []int64, l, r int, mod int64) int64 {
	if l > r {
		return 0
	}
	res := (pref[r+1] - pref[l]*pow[r-l+1]) % mod
	if res < 0 {
		res += mod
	}
	return res
}
