package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// buildPalindrome constructs the lexicographically smallest palindrome
// using the provided letter counts.
func buildPalindrome(freq []int) string {
	left := make([]byte, 0)
	var mid byte
	for i := 0; i < 26; i++ {
		for j := 0; j < freq[i]/2; j++ {
			left = append(left, byte('a'+i))
		}
		freq[i] %= 2
	}
	for i := 0; i < 26; i++ {
		if freq[i] == 1 {
			mid = byte('a' + i)
			break
		}
	}
	right := make([]byte, len(left))
	for i := range left {
		right[i] = left[len(left)-1-i]
	}
	if mid != 0 {
		return string(left) + string(mid) + string(right)
	}
	return string(left) + string(right)
}

// solve returns a string t_max that approximates the minimal possible
// lexicographical maximum of a permutation of s and its reverse.
func solve(s string) string {
	n := len(s)
	freq := make([]int, 26)
	for i := 0; i < n; i++ {
		freq[s[i]-'a']++
	}
	odd := 0
	for i := 0; i < 26; i++ {
		if freq[i]%2 == 1 {
			odd++
		}
	}
	if odd <= 1 {
		f := make([]int, 26)
		copy(f, freq)
		return buildPalindrome(f)
	}

	bytes := []byte(s)
	sort.Slice(bytes, func(i, j int) bool { return bytes[i] < bytes[j] })
	t := make([]byte, n)
	i, j := 0, n-1
	for k := 0; k+1 < n; k += 2 {
		t[i] = bytes[k]
		t[j] = bytes[k+1]
		i++
		j--
	}
	if n%2 == 1 {
		t[i] = bytes[n-1]
	}
	rev := make([]byte, n)
	for x := 0; x < n; x++ {
		rev[x] = t[n-1-x]
	}
	ta := string(t)
	ra := string(rev)
	if ta > ra {
		return ta
	}
	return ra
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, solve(s))
	}
}
