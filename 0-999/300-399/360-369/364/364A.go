package main

import (
	"fmt"
)

func main() {
	var a int64
	var s string
	if _, err := fmt.Scan(&a); err != nil {
		return
	}
	if _, err := fmt.Scan(&s); err != nil {
		return
	}
	n := len(s)
	// prefix sums of digits
	ps := make([]int64, n+1)
	for i := 0; i < n; i++ {
		ps[i+1] = ps[i] + int64(s[i]-'0')
	}
	maxSum := ps[n]
	// freq[sum] = number of subarrays with given sum
	freq := make([]int64, maxSum+1)
	// compute all subarray sums
	for i := 0; i < n; i++ {
		for j := i + 1; j <= n; j++ {
			sum := ps[j] - ps[i]
			freq[sum]++
		}
	}
	total := int64(n) * int64(n+1) / 2
	var ans int64
	if a == 0 {
		f0 := freq[0]
		// count pairs where at least one subarray sum is zero
		ans = 2*f0*total - f0*f0
	} else {
		// iterate over divisors of a
		for d := int64(1); d*d <= a; d++ {
			if a%d != 0 {
				continue
			}
			p := d
			q := a / d
			if p <= maxSum && q <= maxSum {
				if p == q {
					ans += freq[p] * freq[q]
				} else {
					ans += freq[p] * freq[q] * 2
				}
			}
		}
	}
	fmt.Println(ans)
}
