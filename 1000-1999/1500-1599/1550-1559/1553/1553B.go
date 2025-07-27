package main

import (
	"bufio"
	"fmt"
	"os"
)

func canForm(s, t string) bool {
	n := len(s)
	m := len(t)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			k := j - i + 1
			if k > m {
				break
			}
			ok := true
			for p := 0; p < k; p++ {
				if s[i+p] != t[p] {
					ok = false
					break
				}
			}
			if !ok {
				break
			}
			rem := m - k
			if rem == 0 {
				return true
			}
			if j-rem < 0 {
				continue
			}
			ok2 := true
			for p := 0; p < rem; p++ {
				if s[j-1-p] != t[k+p] {
					ok2 = false
					break
				}
			}
			if ok2 {
				return true
			}
		}
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		if canForm(s, t) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
