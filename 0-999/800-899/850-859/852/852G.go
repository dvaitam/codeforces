package main

import (
	"bufio"
	"fmt"
	"os"
)

func expand(pattern []byte, pos int, cur []byte, set map[string]struct{}) {
	if pos == len(pattern) {
		set[string(cur)] = struct{}{}
		return
	}
	ch := pattern[pos]
	if ch == '?' {
		// empty replacement
		expand(pattern, pos+1, cur, set)
		// replace with letters 'a' to 'e'
		for c := byte('a'); c <= 'e'; c++ {
			cur = append(cur, c)
			expand(pattern, pos+1, cur, set)
			cur = cur[:len(cur)-1]
		}
	} else {
		cur = append(cur, ch)
		expand(pattern, pos+1, cur, set)
		cur = cur[:len(cur)-1]
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	words := make(map[string]int, n)
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(in, &w)
		words[w]++
	}

	for i := 0; i < m; i++ {
		var p string
		fmt.Fscan(in, &p)
		set := make(map[string]struct{})
		expand([]byte(p), 0, []byte{}, set)
		total := 0
		for s := range set {
			if cnt, ok := words[s]; ok {
				total += cnt
			}
		}
		fmt.Fprintln(out, total)
	}
}
