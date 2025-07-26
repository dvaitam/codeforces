package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		cnt := make([]int, 26)
		for i := 0; i < 2*n; i++ {
			var s string
			fmt.Fscan(in, &s)
			for j := 0; j < len(s); j++ {
				cnt[s[j]-'a'] ^= 1
			}
		}
		var final string
		fmt.Fscan(in, &final)
		for j := 0; j < len(final); j++ {
			cnt[final[j]-'a'] ^= 1
		}
		for i := 0; i < 26; i++ {
			if cnt[i] == 1 {
				fmt.Fprintf(out, "%c\n", byte('a'+i))
				break
			}
		}
	}
}
