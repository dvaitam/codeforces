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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		var aStr, bStr string
		fmt.Fscan(reader, &aStr, &bStr)
		s := make([][26]int, n+1)
		for i := 0; i < n; i++ {
			s[i+1] = s[i]
			s[i+1][aStr[i]-'a']++
			s[i+1][bStr[i]-'a']--
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l--
			ans := 0
			for i := 0; i < 26; i++ {
				diff := s[r][i] - s[l][i]
				if diff > 0 {
					ans += diff
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
