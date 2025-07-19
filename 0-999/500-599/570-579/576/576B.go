package main

import (
	"bufio"
	"fmt"
	"os"
)

var rdr = bufio.NewReader(os.Stdin)
var wrtr = bufio.NewWriter(os.Stdout)

func readInt() int {
	var x int
	var sign = 1
	c, err := rdr.ReadByte()
	for err == nil && (c < '0' || c > '9') && c != '-' {
		c, err = rdr.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = rdr.ReadByte()
	}
	for err == nil && c >= '0' && c <= '9' {
		x = x*10 + int(c-'0')
		c, err = rdr.ReadByte()
	}
	return x * sign
}

func main() {
	defer wrtr.Flush()
	n := readInt()
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = readInt()
	}
	// check for fixed point
	for i := 1; i <= n; i++ {
		if p[i] == i {
			fmt.Fprintln(wrtr, "YES")
			for j := 1; j <= n; j++ {
				if j != i {
					fmt.Fprintf(wrtr, "%d %d\n", i, j)
				}
			}
			return
		}
	}
	vis := make([]bool, n+1)
	cpar := make([]bool, n+1)
	var a, b int
	good := false
	for i := 1; i <= n; i++ {
		if !vis[i] {
			x := i
			vis[x] = true
			length := 1
			for p[x] != i {
				x = p[x]
				vis[x] = true
				cpar[x] = (length%2 == 1)
				length++
			}
			if length%2 == 1 {
				fmt.Fprintln(wrtr, "NO")
				return
			}
			if length == 2 {
				good = true
				a = i
				b = p[i]
			}
		}
	}
	if !good {
		fmt.Fprintln(wrtr, "NO")
		return
	}
	fmt.Fprintln(wrtr, "YES")
	// root edge
	fmt.Fprintf(wrtr, "%d %d\n", a, b)
	// other edges
	for i := 1; i <= n; i++ {
		if i == a || i == b {
			continue
		}
		if cpar[i] {
			fmt.Fprintf(wrtr, "%d %d\n", a, i)
		} else {
			fmt.Fprintf(wrtr, "%d %d\n", b, i)
		}
	}
}
