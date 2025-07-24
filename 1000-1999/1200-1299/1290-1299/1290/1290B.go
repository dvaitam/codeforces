package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	n := len(s)
	prefix := make([][26]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i]
		prefix[i+1][s[i]-'a']++
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		if l == r {
			fmt.Fprintln(writer, "Yes")
			continue
		}
		if s[l-1] != s[r-1] {
			fmt.Fprintln(writer, "Yes")
			continue
		}
		distinct := 0
		for c := 0; c < 26; c++ {
			if prefix[r][c]-prefix[l-1][c] > 0 {
				distinct++
			}
		}
		if distinct >= 3 {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
