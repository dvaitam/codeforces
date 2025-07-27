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
		a := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			a[i] = []byte(s)
		}
		b := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(in, &s)
			b[i] = []byte(s)
		}

		row := make([]int, n)
		col := make([]int, n)
		for i := 0; i < n; i++ {
			if a[i][0] != b[i][0] {
				row[i] = 1
			}
		}
		for j := 0; j < n; j++ {
			if ((a[0][j] - '0') ^ byte(row[0])) != (b[0][j] - '0') {
				col[j] = 1
			}
		}
		ok := true
	outer:
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				v := int(a[i][j]-'0') ^ row[i] ^ col[j]
				if v != int(b[i][j]-'0') {
					ok = false
					break outer
				}
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
