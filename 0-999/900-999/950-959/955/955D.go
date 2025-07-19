package main

import (
	"bufio"
	"fmt"
	"os"
)

// matchPos finds, for each prefix length l of pat (0<=l<=len(pat)),
// the earliest index in s where a match of length l ends under window k.
func matchPos(s, pat string, k int) []int {
	n, m := len(s), len(pat)
	ans := make([]int, m+1)
	for i := range ans {
		ans[i] = n
	}
	// build prefix function pi for pat
	pi := make([]int, m)
	for i := 1; i < m; i++ {
		j := pi[i-1]
		for j > 0 && pat[j] != pat[i] {
			j = pi[j-1]
		}
		if pat[j] == pat[i] {
			j++
		}
		pi[i] = j
	}
	// build kmpArr: kmpArr[0]=0, kmpArr[i]=pi[i-1] for i>=1
	kmpArr := make([]int, m+1)
	for i := 0; i < m; i++ {
		kmpArr[i+1] = pi[i]
	}
	cur := 0
	for i := 0; i < n; i++ {
		if cur == m {
			cur = kmpArr[cur]
		}
		for cur > 0 && pat[cur] != s[i] {
			cur = kmpArr[cur]
		}
		if pat[cur] == s[i] {
			cur++
		} else {
			cur = 0
		}
		if i+1 < k {
			if cur == m && ans[cur] > k-1 {
				ans[cur] = k - 1
			}
			continue
		}
		tcur := cur
		for tcur > 0 && ans[tcur] > i {
			ans[tcur] = i
			tcur = kmpArr[tcur]
		}
		if ans[tcur] > i {
			ans[tcur] = i
		}
	}
	return ans
}

func reverse(s string) string {
	b := []byte(s)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	var s, t string
	fmt.Fscan(reader, &s, &t)
	// forward match positions for t prefix of length k
	pref := t
	if k < len(t) {
		pref = t[:k]
	}
	forw := matchPos(s, pref, k)
	// reverse strings
	rs := reverse(s)
	rt := reverse(t)
	sufPat := rt
	if k < len(rt) {
		sufPat = rt[:k]
	}
	rev := matchPos(rs, sufPat, k)
	// search for valid split
	for i := 0; i < len(forw); i++ {
		j := len(t) - i
		if j < 0 || j >= len(rev) {
			continue
		}
		if forw[i]+rev[j]+1 < len(s) {
			fmt.Fprintln(writer, "Yes")
			// compute 1-based positions
			start1 := forw[i] - k + 2
			start2 := len(s) - rev[j]
			fmt.Fprintf(writer, "%d %d\n", start1, start2)
			return
		}
	}
	fmt.Fprintln(writer, "No")
}
