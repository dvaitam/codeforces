package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, q int
	fmt.Fscan(in, &n, &q)

	var trans [6][]string
	for i := 0; i < q; i++ {
		var a, b string
		fmt.Fscan(in, &a, &b)
		trans[b[0]-'a'] = append(trans[b[0]-'a'], a)
	}

	cur := map[string]struct{}{"a": {}}
	for length := 1; length < n; length++ {
		next := make(map[string]struct{})
		for s := range cur {
			idx := s[0] - 'a'
			for _, p := range trans[idx] {
				next[p+s[1:]] = struct{}{}
			}
		}
		cur = next
		if len(cur) == 0 {
			break
		}
	}
	fmt.Println(len(cur))
}
