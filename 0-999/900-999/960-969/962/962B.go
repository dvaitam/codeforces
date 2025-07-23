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

	var n, a, b int
	if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)
	prev := 0 // 0 none, 1 programmer(A), 2 athlete(B)
	ans := 0
	for i := 0; i < n; i++ {
		if s[i] == '*' {
			prev = 0
			continue
		}
		if prev == 1 {
			if b > 0 {
				b--
				ans++
				prev = 2
			} else {
				prev = 0
			}
		} else if prev == 2 {
			if a > 0 {
				a--
				ans++
				prev = 1
			} else {
				prev = 0
			}
		} else {
			if a >= b {
				if a > 0 {
					a--
					ans++
					prev = 1
				} else if b > 0 {
					b--
					ans++
					prev = 2
				} else {
					prev = 0
				}
			} else {
				if b > 0 {
					b--
					ans++
					prev = 2
				} else if a > 0 {
					a--
					ans++
					prev = 1
				} else {
					prev = 0
				}
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
