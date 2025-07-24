package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)
		m := n / k
		freq := make([]int, 26)
		for i := 0; i < n; i++ {
			freq[int(s[i]-'a')]++
		}
		size := make([]int, k)
		mex := make([]byte, k)
		for letter := 0; letter < 26; letter++ {
			active := make([]int, 0, k)
			for i := 0; i < k; i++ {
				if mex[i] == 0 && size[i] < m {
					active = append(active, i)
				}
			}
			use := freq[letter]
			if use > len(active) {
				use = len(active)
			}
			for idx := 0; idx < use; idx++ {
				i := active[idx]
				size[i]++
				freq[letter]--
				if size[i] == m {
					mex[i] = byte('a' + letter + 1)
				}
			}
			for idx := use; idx < len(active); idx++ {
				i := active[idx]
				mex[i] = byte('a' + letter)
			}
		}
		for i := 0; i < k; i++ {
			if mex[i] == 0 {
				mex[i] = 'z'
			}
		}
		sort.Slice(mex, func(i, j int) bool { return mex[i] > mex[j] })
		fmt.Fprintln(writer, string(mex))
	}
}
