package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	var ss, sp string
	fmt.Fscan(reader, &ss)
	fmt.Fscan(reader, &sp)
	lt := n
	ls := len(ss)
	lp := len(sp)
	st := make([]byte, n)
	for i := range st {
		st[i] = '!'
	}
	tt := 0
	for i := 0; i < lp; i++ {
		if sp[i] == '1' {
			tt++
			for j := 0; j < ls; j++ {
				idx := i + j
				if st[idx] != '!' {
					if st[idx] != ss[j] {
						fmt.Println("No solution")
						return
					}
				} else {
					st[idx] = ss[j]
				}
			}
		}
	}
	cheak := func(s []byte) bool {
		tot := 0
		for i := 0; i+ls <= lt; i++ {
			f := true
			for j := 0; j < ls; j++ {
				if s[i+j] != ss[j] {
					f = false
					break
				}
			}
			if f {
				tot++
			}
		}
		return tot == tt
	}
	s := make([]byte, n)
	for ch := byte('a'); ch < byte('a')+byte(k); ch++ {
		for j := 0; j < lt; j++ {
			if st[j] == '!' {
				s[j] = ch
			} else {
				s[j] = st[j]
			}
		}
		if cheak(s) {
			fmt.Println(string(s))
			return
		}
	}
	fmt.Println("No solution")
}
