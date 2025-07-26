package main

import (
	"bufio"
	"fmt"
	"os"
)

func isGood(s string) bool {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == s[i+1] {
			return false
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)

		if isGood(s) {
			fmt.Fprintln(writer, "YES")
			continue
		}
		if !isGood(t) {
			fmt.Fprintln(writer, "NO")
			continue
		}

		has00, has11 := false, false
		for i := 0; i+1 < len(s); i++ {
			if s[i] == '0' && s[i+1] == '0' {
				has00 = true
			}
			if s[i] == '1' && s[i+1] == '1' {
				has11 = true
			}
		}

		if has00 && has11 {
			fmt.Fprintln(writer, "NO")
			continue
		}

		first, last := t[0], t[len(t)-1]
		if has00 {
			if first == '1' && last == '1' {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		} else if has11 {
			if first == '0' && last == '0' {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
