package main

import (
	"bufio"
	"fmt"
	"os"
)

const MAXN = 1000000

var spf [MAXN + 1]int

func init() {
	for i := 2; i <= MAXN; i++ {
		if spf[i] == 0 {
			for j := i; j <= MAXN; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		if n == 1 || m == 1 {
			fmt.Fprintln(writer, "YES")
			continue
		}
		d := spf[n]
		if d == 0 { // n == 1, already handled
			d = n
		}
		if m >= d {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
