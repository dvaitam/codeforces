package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt (contest 1619).
// Given integers a and s we reconstruct b so that performing Tanya's "wrong"
// digit-wise addition of a and b yields s. We process digits from right to left.
// If the current digit of s is smaller than that of a we borrow the next digit
// of s to form a two-digit number. If borrowing fails or digits remain in a
// after s is exhausted, then b doesn't exist.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var aStr, sStr string
		fmt.Fscan(in, &aStr, &sStr)
		res := solve(aStr, sStr)
		fmt.Fprintln(out, res)
	}
}

func solve(aStr, sStr string) string {
	i := len(aStr) - 1
	j := len(sStr) - 1
	ans := make([]byte, 0, len(sStr))
	for i >= 0 {
		if j < 0 {
			return "-1"
		}
		aDigit := aStr[i] - '0'
		sDigit := sStr[j] - '0'
		if sDigit >= aDigit {
			ans = append(ans, '0'+(sDigit-aDigit))
			j--
		} else {
			if j == 0 {
				return "-1"
			}
			twoDigit := (sStr[j-1]-'0')*10 + sDigit
			diff := twoDigit - aDigit
			if diff < 0 || diff > 9 {
				return "-1"
			}
			ans = append(ans, '0'+diff)
			j -= 2
		}
		i--
	}
	for j >= 0 {
		ans = append(ans, sStr[j])
		j--
	}
	// remove leading zeros
	for len(ans) > 1 && ans[len(ans)-1] == '0' {
		ans = ans[:len(ans)-1]
	}
	// reverse ans
	for l, r := 0, len(ans)-1; l < r; l, r = l+1, r-1 {
		ans[l], ans[r] = ans[r], ans[l]
	}
	return string(ans)
}
