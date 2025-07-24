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

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var s, t string
		fmt.Fscan(in, &s)
		fmt.Fscan(in, &t)

		pos := make([][]int, 26)
		for i, ch := range s {
			pos[ch-'a'] = append(pos[ch-'a'], i)
		}

		possible := true
		// quick check: each char in t must appear in s
		for _, ch := range t {
			if len(pos[ch-'a']) == 0 {
				possible = false
				break
			}
		}
		if !possible {
			fmt.Fprintln(out, -1)
			continue
		}

		operations := 1
		idx := 0
		i := 0
		for i < len(t) {
			ch := t[i]
			arr := pos[ch-'a']
			j := sort.Search(len(arr), func(k int) bool { return arr[k] >= idx })
			if j == len(arr) {
				operations++
				idx = 0
			} else {
				idx = arr[j] + 1
				i++
			}
		}
		fmt.Fprintln(out, operations)
	}
}
