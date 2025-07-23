package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t string
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	n := len(t)
	for k := 1; k < n; k++ {
		if (n+k)%2 != 0 {
			continue
		}
		L := (n + k) / 2
		if k >= L || L > n {
			continue
		}
		s := t[:L]
		if t[n-L:] != s {
			continue
		}
		if L < len(t) {
			if t[L:] != s[k:] {
				continue
			}
		} else {
			if k != L {
				continue
			}
		}
		fmt.Println("YES")
		fmt.Println(s)
		return
	}
	fmt.Println("NO")
}
