package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		var s string
		fmt.Fscan(reader, &s)
		mp := make(map[int]byte)
		ok := true
		for i := 0; i < n && ok; i++ {
			ch := s[i]
			if v, exists := mp[a[i]]; exists {
				if v != ch {
					ok = false
				}
			} else {
				mp[a[i]] = ch
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
