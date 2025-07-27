package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var enemy, gregor string
		fmt.Fscan(in, &enemy)
		fmt.Fscan(in, &gregor)
		e := []byte(enemy)
		g := []byte(gregor)
		ans := 0
		for i := 0; i < n; i++ {
			if g[i] == '1' {
				if e[i] == '0' {
					ans++
					// occupy this square so other pawns cannot use it
					e[i] = '2'
				} else if i > 0 && e[i-1] == '1' {
					ans++
					e[i-1] = '2'
				} else if i+1 < n && e[i+1] == '1' {
					ans++
					e[i+1] = '2'
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
