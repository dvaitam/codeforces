package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxBeauty(s string, n int) int {
	freq := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		freq[s[i]]++
	}
	maxFreq := 0
	for _, v := range freq {
		if v > maxFreq {
			maxFreq = v
		}
	}
	l := len(s)
	if maxFreq == l && n == 1 {
		return l - 1
	}
	if maxFreq+n > l {
		return l
	}
	return maxFreq + n
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s1, s2, s3 string
	fmt.Fscan(in, &s1, &s2, &s3)

	b1 := maxBeauty(s1, n)
	b2 := maxBeauty(s2, n)
	b3 := maxBeauty(s3, n)

	if b1 > b2 && b1 > b3 {
		fmt.Fprintln(out, "Kuro")
	} else if b2 > b1 && b2 > b3 {
		fmt.Fprintln(out, "Shiro")
	} else if b3 > b1 && b3 > b2 {
		fmt.Fprintln(out, "Katie")
	} else {
		fmt.Fprintln(out, "Draw")
	}
}
