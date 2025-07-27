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
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		if n%k != 0 {
			fmt.Fprintln(out, -1)
			continue
		}
		freq := make([]int, 26)
		for i := 0; i < n; i++ {
			freq[int(s[i]-'a')]++
		}
		good := true
		for i := 0; i < 26; i++ {
			if freq[i]%k != 0 {
				good = false
				break
			}
		}
		if good {
			fmt.Fprintln(out, s)
			continue
		}
		ans := ""
		found := false
		// try to change position from end
		for i := n - 1; i >= 0 && !found; i-- {
			cur := int(s[i] - 'a')
			freq[cur]--
			for c := cur + 1; c < 26 && !found; c++ {
				freq[c]++
				rem := n - i - 1
				needSum := 0
				for j := 0; j < 26; j++ {
					needSum += (k - freq[j]%k) % k
				}
				if needSum <= rem && (rem-needSum)%k == 0 {
					// construct result
					prefix := s[:i] + string('a'+byte(c))
					rest := make([]byte, 0, rem)
					for t := 0; t < rem-needSum; t++ {
						rest = append(rest, 'a')
					}
					for j := 0; j < 26; j++ {
						cnt := (k - freq[j]%k) % k
						for t := 0; t < cnt; t++ {
							rest = append(rest, byte('a'+j))
						}
					}
					// sort rest since appended in order ensures lexicographic
					// since we appended 'a's first then increasing letters
					ans = prefix + string(rest)
					found = true
				}
				freq[c]--
			}
		}
		if found {
			fmt.Fprintln(out, ans)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
