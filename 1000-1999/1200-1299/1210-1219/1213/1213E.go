package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var s, t string
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	fmt.Fscan(in, &s, &t)
	// Case 1: swapped characters
	if s[0] != s[1] && t[0] != t[1] && s[0] == t[1] && s[1] == t[0] {
		fmt.Println("YES")
		// build sequence [s[0], missing, s[1]]
		var mid byte
		for _, c := range []byte{'a', 'b', 'c'} {
			if c != s[0] && c != s[1] {
				mid = c
			}
		}
		out := []byte{s[0], mid, s[1]}
		var sb strings.Builder
		for i := 0; i < 3; i++ {
			for j := 0; j < n; j++ {
				sb.WriteByte(out[i])
			}
		}
		fmt.Println(sb.String())
		return
	}
	// Case 2: same first char
	if s[0] == t[0] && s[1] != t[1] && s[0] != s[1] && t[0] != t[1] {
		fmt.Println("YES")
		var out []byte
		for _, c := range []byte{'a', 'b', 'c'} {
			if c != s[0] {
				out = append(out, c)
			}
		}
		out = append(out, s[0])
		var sb strings.Builder
		for i := 0; i < 3; i++ {
			for j := 0; j < n; j++ {
				sb.WriteByte(out[i])
			}
		}
		fmt.Println(sb.String())
		return
	}
	// Case 3: same second char
	if s[1] == t[1] && s[0] != t[0] && s[0] != s[1] && t[0] != t[1] {
		fmt.Println("YES")
		var out []byte
		out = append(out, s[1])
		for _, c := range []byte{'a', 'b', 'c'} {
			if c != s[1] {
				out = append(out, c)
			}
		}
		var sb strings.Builder
		for i := 0; i < 3; i++ {
			for j := 0; j < n; j++ {
				sb.WriteByte(out[i])
			}
		}
		fmt.Println(sb.String())
		return
	}
	// General case: try all permutations
	perms := [][]byte{{'a', 'b', 'c'}, {'a', 'c', 'b'}, {'b', 'a', 'c'}, {'b', 'c', 'a'}, {'c', 'a', 'b'}, {'c', 'b', 'a'}}
	for _, p := range perms {
		// build circular pp
		pp := []byte{p[0], p[1], p[2], p[0]}
		ok := true
		for _, q := range [][]byte{[]byte(s), []byte(t)} {
			for i := 0; i < 3; i++ {
				if q[0] == pp[i] && q[1] == pp[i+1] {
					ok = false
				}
			}
		}
		if ok {
			fmt.Println("YES")
			var sb strings.Builder
			for i := 0; i < n; i++ {
				sb.Write(p)
			}
			fmt.Println(sb.String())
			return
		}
	}
	fmt.Println("NO")
}
