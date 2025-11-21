package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		cnt := make([]int, 26)
		for _, ch := range s {
			cnt[ch-'a']++
		}
		type pair struct {
			c byte
			f int
		}
		pairs := make([]pair, 0, 26)
		for i := 0; i < 26; i++ {
			if cnt[i] > 0 {
				pairs = append(pairs, pair{c: byte('a' + i), f: cnt[i]})
			}
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].f == pairs[j].f {
				return pairs[i].c < pairs[j].c
			}
			return pairs[i].f > pairs[j].f
		})

		if len(pairs) == 1 {
			fmt.Fprintln(out, s)
			continue
		}

		// place characters in alternating fashion to maximize alternations
		res := make([]byte, n)
		pos := 0
		for idx := 0; idx < len(pairs); idx++ {
			for pairs[idx].f > 0 && pos < n {
				res[pos] = pairs[idx].c
				pos += 2
				pairs[idx].f--
			}
		}
		pos = 1
		for idx := 0; idx < len(pairs); idx++ {
			for pairs[idx].f > 0 && pos < n {
				res[pos] = pairs[idx].c
				pos += 2
				pairs[idx].f--
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
