package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	words := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		words[s] = struct{}{}
	}
	var q int
	fmt.Fscan(in, &q)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; q > 0; q-- {
		var t string
		fmt.Fscan(in, &t)
		seen := make(map[string]struct{})
		cnt := 0
		for i := 0; i < len(t); i++ {
			cand := t[:i] + t[i+1:]
			if _, ok := words[cand]; ok {
				if _, used := seen[cand]; !used {
					cnt++
					seen[cand] = struct{}{}
				}
			}
		}
		fmt.Fprintln(out, cnt)
	}
}
