package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// largestBeautiful returns the largest beautiful number less than s.
// A beautiful number has even length and its digits can be permuted into a palindrome,
// which for even length means every digit occurs an even number of times.
func largestBeautiful(s string) string {
	n := len(s)
	prefix := make([][10]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i]
		prefix[i+1][s[i]-'0']++
	}

	// try to keep prefix equal to s as long as possible
	for i := n - 1; i >= 0; i-- {
		cur := int(s[i] - '0')
		for d := cur - 1; d >= 0; d-- {
			if i == 0 && d == 0 { // leading zero not allowed
				continue
			}
			cnt := prefix[i]
			cnt[d]++
			remaining := n - i - 1
			odd := []int{}
			for k := 0; k < 10; k++ {
				if cnt[k]%2 == 1 {
					odd = append(odd, k)
				}
			}
			if remaining < len(odd) || (remaining-len(odd))%2 != 0 {
				continue
			}
			// construct the maximal suffix
			rest := remaining - len(odd)
			digits := make([]int, 0, remaining)
			digits = append(digits, odd...)
			for j := 0; j < rest; j++ {
				digits = append(digits, 9)
			}
			sort.Slice(digits, func(a, b int) bool { return digits[a] > digits[b] })
			ans := make([]byte, 0, n)
			ans = append(ans, s[:i]...)
			ans = append(ans, byte('0'+d))
			for _, x := range digits {
				ans = append(ans, byte('0'+x))
			}
			return string(ans)
		}
	}

	// If no number of the same length works, answer has length n-2 consisting of all '9'
	res := make([]byte, n-2)
	for i := range res {
		res[i] = '9'
	}
	return string(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, largestBeautiful(s))
	}
}
