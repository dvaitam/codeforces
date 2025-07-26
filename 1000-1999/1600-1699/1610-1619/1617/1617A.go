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
		var s, perm string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &perm)

		cnt := make([]int, 26)
		for i := 0; i < len(s); i++ {
			cnt[s[i]-'a']++
		}

		if perm == "abc" && cnt[0] > 0 && cnt[1] > 0 && cnt[2] > 0 {
			for i := 0; i < cnt[0]; i++ {
				writer.WriteByte('a')
			}
			for i := 0; i < cnt[2]; i++ {
				writer.WriteByte('c')
			}
			for i := 0; i < cnt[1]; i++ {
				writer.WriteByte('b')
			}
			for ch := byte('d'); ch <= 'z'; ch++ {
				for i := 0; i < cnt[ch-'a']; i++ {
					writer.WriteByte(ch)
				}
			}
			writer.WriteByte('\n')
		} else {
			for ch := byte('a'); ch <= 'z'; ch++ {
				for i := 0; i < cnt[ch-'a']; i++ {
					writer.WriteByte(ch)
				}
			}
			writer.WriteByte('\n')
		}
	}
}
