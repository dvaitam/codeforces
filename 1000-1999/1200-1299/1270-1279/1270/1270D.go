package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Interactive solution for Codeforces problem 1270D - Strange Device.
// The program queries the first k+1 indices, omitting each index once,
// and counts how many times the smallest value among answers occurs.
// That count equals k+1 - m, so the program outputs m.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	type pair struct{ val, pos int }
	ans := make([]pair, k+1)

	for i := 0; i <= k; i++ {
		fmt.Fprint(out, "?")
		for j := 0; j <= k; j++ {
			if j == i {
				continue
			}
			fmt.Fprintf(out, " %d", j+1)
		}
		fmt.Fprintln(out)
		out.Flush()

		var pos, val int
		if _, err := fmt.Fscan(in, &pos, &val); err != nil {
			return
		}
		ans[i] = pair{val: val, pos: pos}
	}

	sort.Slice(ans, func(i, j int) bool { return ans[i].val < ans[j].val })

	smallest := ans[0].val
	cnt := 0
	for _, p := range ans {
		if p.val == smallest {
			cnt++
		}
	}

	m := k + 1 - cnt
	fmt.Fprintf(out, "! %d\n", m)
	out.Flush()
}
