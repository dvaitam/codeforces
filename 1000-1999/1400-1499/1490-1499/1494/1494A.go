package main

import (
	"bufio"
	"fmt"
	"os"
)

func isValid(s []byte) bool {
	bal := 0
	for _, c := range s {
		if c == '(' {
			bal++
		} else {
			bal--
			if bal < 0 {
				return false
			}
		}
	}
	return bal == 0
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a string
		fmt.Fscan(in, &a)
		ok := false
		for mask := 0; mask < 8 && !ok; mask++ {
			mp := map[byte]byte{'A': ')', 'B': ')', 'C': ')'}
			if mask&1 != 0 {
				mp['A'] = '('
			}
			if mask&2 != 0 {
				mp['B'] = '('
			}
			if mask&4 != 0 {
				mp['C'] = '('
			}
			b := make([]byte, len(a))
			for i := 0; i < len(a); i++ {
				b[i] = mp[a[i]]
			}
			if isValid(b) {
				ok = true
			}
		}
		if ok {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}
}
