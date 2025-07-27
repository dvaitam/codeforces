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
		var n int
		fmt.Fscan(reader, &n)
		var s1, s2 string
		fmt.Fscan(reader, &s1)
		fmt.Fscan(reader, &s2)

		ans := 0
		i := 0
		for i < n {
			if s1[i] != s2[i] {
				ans += 2
				i++
				continue
			}
			if s1[i] == '0' {
				if i+1 < n && s1[i+1] == '1' && s2[i+1] == '1' {
					ans += 2
					i += 2
				} else {
					ans++
					i++
				}
			} else { // both '1'
				if i+1 < n && s1[i+1] == '0' && s2[i+1] == '0' {
					ans += 2
					i += 2
				} else {
					i++
				}
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
