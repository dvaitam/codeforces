package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int = 1000000007

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	cur := 0
	ans := 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == 'b' {
			cur++
			if cur >= mod {
				cur -= mod
			}
		} else { // s[i] == 'a'
			ans += cur
			if ans >= mod {
				ans -= mod
			}
			cur <<= 1
			if cur >= mod {
				cur %= mod
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
