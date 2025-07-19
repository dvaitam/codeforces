package main

import (
	"bufio"
	"fmt"
	"os"
)

// ok checks if s[l..r] is a palindrome
func ok(s []byte, l, r int) bool {
	for l <= r && s[l] == s[r] {
		l++
		r--
	}
	return l > r
}

// solve returns the longest palindrome by keeping matching prefix and suffix and
// adding the best palindrome from the middle.
func solve(s string) string {
	b := []byte(s)
	l, r := 0, len(b)-1
	// match prefix and suffix
	for l < r && b[l] == b[r] {
		l++
		r--
	}
	// if the whole string is palindrome
	if l >= r {
		return s
	}
	// find longest prefix palindrome in middle
	r2 := l
	for i := r; i >= l; i-- {
		if ok(b, l, i) {
			r2 = i
			break
		}
	}
	// find longest suffix palindrome in middle
	l2 := r
	for i := l; i <= r; i++ {
		if ok(b, i, r) {
			l2 = i
			break
		}
	}
	// choose the longer one
	var mid []byte
	if r2-l+1 > r-l2+1 {
		mid = b[l : r2+1]
	} else {
		mid = b[l2 : r+1]
	}
	// build result: prefix + mid + suffix
	res := make([]byte, 0, len(b))
	res = append(res, b[:l]...)
	res = append(res, mid...)
	res = append(res, b[r+1:]...)
	return string(res)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for tc := 0; tc < t; tc++ {
		var s string
		fmt.Fscan(reader, &s)
		ans := solve(s)
		fmt.Fprintln(writer, ans)
	}
}
