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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		solve(in, out)
	}
}

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)

	var r1, c1, r2, c2, r3, c3 int
	fmt.Fscan(in, &r1, &c1, &r2, &c2, &r3, &c3)
	var x, y int
	fmt.Fscan(in, &x, &y)

	X := r1
	if r1 == r2 || r1 == r3 {
		X = r1
	} else {
		X = r2
	}
	Y := c1
	if c1 == c2 || c1 == c3 {
		Y = c1
	} else {
		Y = c2
	}

	if (X == 1 && (Y == 1 || Y == n)) || (X == n && (Y == 1 || Y == n)) {
		if x == X || y == Y {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
		return
	}

	if x%2 == X%2 || y%2 == Y%2 {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
