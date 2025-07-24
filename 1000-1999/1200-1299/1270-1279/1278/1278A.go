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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var p, h string
		fmt.Fscan(in, &p)
		fmt.Fscan(in, &h)
		if len(p) > len(h) {
			fmt.Fprintln(out, "NO")
			continue
		}
		var freqP [26]int
		for i := 0; i < len(p); i++ {
			freqP[p[i]-'a']++
		}
		var freqH [26]int
		lenP := len(p)
		for i := 0; i < lenP; i++ {
			freqH[h[i]-'a']++
		}
		found := func() bool {
			for i := 0; i < 26; i++ {
				if freqP[i] != freqH[i] {
					return false
				}
			}
			return true
		}
		ok := found()
		for i := lenP; !ok && i < len(h); i++ {
			freqH[h[i]-'a']++
			freqH[h[i-lenP]-'a']--
			if found() {
				ok = true
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
