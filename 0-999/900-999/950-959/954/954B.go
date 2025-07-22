package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var s string
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	ans := n
	for l := 1; l*2 <= n; l++ {
		if s[:l] == s[l:2*l] {
			ops := n - l + 1
			if ops < ans {
				ans = ops
			}
		}
	}
	fmt.Println(ans)
}
