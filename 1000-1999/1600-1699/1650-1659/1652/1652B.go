package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution to problemB.txt for 1652B.
// The algorithm repeatedly removes the longest prefix that
// appears elsewhere in the current string. This is equivalent
// to discarding all leading characters that appear more than once.
// We count character frequencies and skip characters from the start
// while their remaining frequency is at least two.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		freq := [26]int{}
		for i := 0; i < len(s); i++ {
			freq[s[i]-'a']++
		}
		idx := 0
		for idx < len(s) && freq[s[idx]-'a'] >= 2 {
			freq[s[idx]-'a']--
			idx++
		}
		fmt.Fprintln(out, s[idx:])
	}
}
