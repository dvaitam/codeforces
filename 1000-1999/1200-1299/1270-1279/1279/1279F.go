package main

import (
	"bufio"
	"fmt"
	"os"
)

// maxCoverage computes the maximum number of target letters (1's in arr)
// that can be covered using at most k disjoint segments of length l.
func maxCoverage(arr []int, l, k int) int {
	n := len(arr)
	if k == 0 || l == 0 || n == 0 {
		return 0
	}
	if l > n {
		l = n
	}
	// prefix sums of arr
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefix[i] = prefix[i-1] + arr[i-1]
	}
	// dp arrays: previous and current
	dpPrev := make([]int, n+1)
	dpCurr := make([]int, n+1)
	for seg := 1; seg <= k; seg++ {
		for i := 1; i <= n; i++ {
			// option1: skip position i
			if dpCurr[i-1] > dpCurr[i] {
				dpCurr[i] = dpCurr[i-1]
			}
			if i >= l {
				val := dpPrev[i-l] + prefix[i] - prefix[i-l]
				if val > dpCurr[i] {
					dpCurr[i] = val
				}
			}
		}
		// prepare for next iteration
		for i := 0; i <= n; i++ {
			dpPrev[i] = dpCurr[i]
			dpCurr[i] = 0
		}
	}
	return dpPrev[n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k, l int
	if _, err := fmt.Fscan(in, &n, &k, &l); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)
	lowerArr := make([]int, n)
	upperArr := make([]int, n)
	var lowerCount, upperCount int
	for i := 0; i < n; i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			lowerArr[i] = 1
			lowerCount++
		} else {
			upperArr[i] = 1
			upperCount++
		}
	}
	if k > n/l {
		k = n / l
	}
	// option1: convert uppercase letters to lowercase
	coverUpper := maxCoverage(upperArr, l, k)
	// option2: convert lowercase letters to uppercase
	coverLower := maxCoverage(lowerArr, l, k)
	res1 := upperCount - coverUpper
	res2 := lowerCount - coverLower
	if res1 < res2 {
		fmt.Println(res1)
	} else {
		fmt.Println(res2)
	}
}
