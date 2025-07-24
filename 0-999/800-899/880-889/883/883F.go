package main

import (
	"bufio"
	"fmt"
	"os"
)

func canonical(s string) string {
	res := make([]byte, 0, len(s))
	cnt := 0
	flush := func() {
		for cnt >= 2 {
			res = append(res, 'u')
			cnt -= 2
		}
		if cnt == 1 {
			res = append(res, 'o')
			cnt = 0
		}
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == 'o' {
			cnt++
		} else if c == 'u' {
			cnt += 2
		} else {
			flush()
			if c == 'h' {
				for len(res) > 0 && res[len(res)-1] == 'k' {
					res = res[:len(res)-1]
				}
				res = append(res, 'h')
			} else {
				res = append(res, c)
			}
		}
	}
	flush()
	return string(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	names := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		names[canonical(s)] = struct{}{}
	}
	fmt.Fprintln(out, len(names))
}
