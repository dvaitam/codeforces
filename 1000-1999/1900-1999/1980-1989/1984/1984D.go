package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
)

func validT(s, t string) bool {
	if t == "a" {
		return false
	}
	n := len(s)
	m := len(t)
	i := 0
	for i < n {
		if i+m <= n && s[i:i+m] == t {
			i += m
		} else if s[i] == 'a' {
			i++
		} else {
			return false
		}
	}
	return true
}

func countT(s string) int {
	n := len(s)
	// build r without 'a'
	rBytes := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		if s[i] != 'a' {
			rBytes = append(rBytes, s[i])
		}
	}
	r := string(rBytes)
	rlen := len(r)

	res := make(map[string]struct{})

	if rlen == 0 {
		for L := 2; L <= n; L++ {
			res[string(bytes.Repeat([]byte{'a'}, L))] = struct{}{}
		}
		return len(res)
	}

	divis := []int{}
	for d := 1; d*d <= rlen; d++ {
		if rlen%d == 0 {
			divis = append(divis, d)
			if d*d != rlen {
				divis = append(divis, rlen/d)
			}
		}
	}
	sort.Ints(divis)

	for _, L := range divis {
		rt := r[:L]
		for start := 0; start < n; start++ {
			if start > 0 && s[start-1] != 'a' {
				continue
			}
			j := 0
			i := start
			first := -1
			for i < n && j < L {
				if s[i] == rt[j] {
					if first == -1 {
						first = i
					}
					j++
					i++
				} else if s[i] == 'a' {
					if first == -1 {
						first = i
					}
					i++
				} else {
					break
				}
			}
			if j == L && first != -1 {
				end := i
				k := 0
				for end+k < n && s[end+k] == 'a' {
					k++
				}
				for t := 0; t <= k; t++ {
					cand := s[start : end+t]
					if validT(s, cand) {
						res[cand] = struct{}{}
					}
				}
			}
		}
	}

	return len(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T++ {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, countT(s))
	}
}
