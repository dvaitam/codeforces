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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		freq := [4]int{}
		ans := len(s) + 1
		l := 0
		for r := 0; r < len(s); r++ {
			freq[s[r]-'0']++
			for freq[1] > 0 && freq[2] > 0 && freq[3] > 0 {
				if r-l+1 < ans {
					ans = r - l + 1
				}
				freq[s[l]-'0']--
				l++
			}
		}
		if ans == len(s)+1 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}
