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
		var s string
		fmt.Fscan(reader, &s)
		allA := true
		for i := 0; i < len(s); i++ {
			if s[i] != 'a' {
				allA = false
				break
			}
		}
		if allA {
			fmt.Fprintln(writer, "NO")
			continue
		}
		left, right := 0, len(s)-1
		ans := 0 // 0: prepend 'a', 1: append 'a'
		for left < right {
			if s[left] == 'a' && s[right] == 'a' {
				left++
				right--
			} else {
				if s[left] == 'a' {
					ans = 0
				} else if s[right] == 'a' {
					ans = 1
				} else {
					ans = 0
				}
				break
			}
		}
		if left > right {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		if ans == 0 {
			// prepend 'a'
			fmt.Fprint(writer, "a")
			fmt.Fprintln(writer, s)
		} else {
			// append 'a'
			fmt.Fprintln(writer, s+"a")
		}
	}
}
