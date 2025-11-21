package main

import (
	"bufio"
	"fmt"
	"os"
)

func canRearrange(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	var freq [26]int
	for i := 0; i < len(s); i++ {
		freq[s[i]-'a']++
		freq[t[i]-'a']--
	}
	for _, v := range freq {
		if v != 0 {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		var s, t string
		fmt.Fscan(in, &s, &t)
		if canRearrange(s, t) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
