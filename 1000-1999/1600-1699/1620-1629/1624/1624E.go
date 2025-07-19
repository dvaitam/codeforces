package main

import (
	"bufio"
	"fmt"
	"os"
)

// info holds the segment boundaries and the source string index
type info struct{ l, r, id int }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var tt int
	fmt.Fscan(in, &tt)
	for tt > 0 {
		tt--
		var n, m int
		fmt.Fscan(in, &n, &m)
		map2 := make(map[string]info)
		map3 := make(map[string]info)
		var s string
		// read source strings and record all substrings of length 2 and 3
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &s)
			for j := 0; j+1 < m; j++ {
				key := s[j : j+2]
				map2[key] = info{j + 1, j + 2, i}
			}
			for j := 0; j+2 < m; j++ {
				key := s[j : j+3]
				map3[key] = info{j + 1, j + 3, i}
			}
		}
		// read target string
		fmt.Fscan(in, &s)
		// dp arrays
		valid := make([]bool, m)
		prev := make([]int, m)
		seg := make([]info, m)
		for i := range prev {
			prev[i] = -1
		}
		// initialize for length 2 and 3 prefixes
		if m >= 2 {
			if v, ok := map2[s[0:2]]; ok {
				valid[1] = true
				seg[1] = v
			}
		}
		if m >= 3 {
			if v, ok := map3[s[0:3]]; ok {
				valid[2] = true
				seg[2] = v
			}
		}
		// dp for positions >= 3
		for i := 3; i < m; i++ {
			key2 := s[i-1 : i+1]
			if v, ok := map2[key2]; ok && valid[i-2] {
				valid[i] = true
				prev[i] = i - 2
				seg[i] = v
			} else {
				key3 := s[i-2 : i+1]
				if v3, ok3 := map3[key3]; ok3 && valid[i-3] {
					valid[i] = true
					prev[i] = i - 3
					seg[i] = v3
				}
			}
		}
		// no solution
		if m == 0 || !valid[m-1] {
			fmt.Fprintln(out, -1)
			continue
		}
		// reconstruct solution
		var res []info
		for i := m - 1; i >= 0; {
			res = append(res, seg[i])
			if prev[i] < 0 {
				break
			}
			i = prev[i]
		}
		// reverse
		for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
			res[i], res[j] = res[j], res[i]
		}
		// output
		fmt.Fprintln(out, len(res))
		for _, v := range res {
			fmt.Fprintln(out, v.l, v.r, v.id)
		}
	}
}
