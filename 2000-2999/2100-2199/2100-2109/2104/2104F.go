package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

var v []int

func dfs(id int, now int, qd bool, s map[string]bool) {
	if id == 0 {
		if now >= 1 && now+2 <= 1000000000 {
			s1 := strconv.Itoa(now)
			s2 := strconv.Itoa(now + 1)
			b := []byte(s1 + s2)
			sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
			gg := string(b)
			if !s[gg] {
				s[gg] = true
				v = append(v, now)
			}
		}
		return
	}
	if qd {
		dfs(id-1, now*10+9, qd, s)
		return
	}
	w := 0
	if now >= 10 {
		w = now % 10
	}
	for i := 0; i < 10; i++ {
		dfs(id-1, now*10+i, qd || (w > i), s)
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	s := make(map[string]bool)
	dfs(9, 0, false, s)
	sort.Ints(v)

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for i := 0; i < q; i++ {
		var r int
		fmt.Fscan(in, &r)
		idx := sort.SearchInts(v, r+1)
		fmt.Fprintln(out, idx)
	}
}