package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solves problem 557E - Ann and Half-Palindrome.
// It outputs the k-th lexicographical substring of s that is a half-palindrome.
func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	var k int
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	n := len(s)

	// dp[i][j] is true if s[i..j] (inclusive) is a half-palindrome.
	dp := make([][]bool, n)
	for i := range dp {
		dp[i] = make([]bool, n)
	}
	for l := 1; l <= n; l++ {
		for i := 0; i+l-1 < n; i++ {
			j := i + l - 1
			if l == 1 {
				dp[i][j] = true
			} else if l == 2 {
				dp[i][j] = s[i] == s[j]
			} else if s[i] == s[j] && (l <= 4 || dp[i+2][j-2]) {
				dp[i][j] = true
			}
		}
	}

	// pre[i][L] — number of half-palindromic substrings starting at i
	// having length at least L.
	pre := make([][]int, n)
	for i := n - 1; i >= 0; i-- {
		pre[i] = make([]int, n-i+2)
		for L := n - i; L >= 1; L-- {
			pre[i][L] = pre[i][L+1]
			if dp[i][i+L-1] {
				pre[i][L]++
			}
		}
	}

	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}
	prefix := make([]byte, 0, n)
	L := 0
	for {
		found := false
		for _, ch := range []byte{'a', 'b'} {
			newIdx := make([]int, 0, len(indices))
			for _, idx := range indices {
				if idx+L < n && s[idx+L] == ch {
					newIdx = append(newIdx, idx)
				}
			}
			if len(newIdx) == 0 {
				continue
			}
			total := 0
			for _, idx := range newIdx {
				total += pre[idx][L+1]
				if total >= k {
					break
				}
			}
			if total >= k {
				prefix = append(prefix, ch)
				indices = newIdx
				L++
				found = true
				break
			} else {
				k -= total
			}
		}
		if !found {
			return
		}
		equal := 0
		for _, idx := range indices {
			if dp[idx][idx+L-1] {
				equal++
			}
		}
		if equal >= k {
			fmt.Println(string(prefix))
			return
		}
		k -= equal
	}
}
