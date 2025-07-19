package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var a, b, c, d int
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &a, &b, &c, &d); err != nil {
		return
	}
	// initialize grid
	m := make([][]byte, 50)
	for i := 0; i < 50; i++ {
		m[i] = make([]byte, 50)
		for j := 0; j < 50; j++ {
			m[i][j] = '#'
		}
	}
	// helper functions
	f := func(x, y int, ch byte) {
		w := 0
		for i := 1; i <= x; i++ {
			row := y - 1 + w%3
			col := i
			m[row][col] = ch
			w++
		}
	}
	ff := func(k, x, y int, ch byte) {
		w := 0
		for i := 1; i <= x; i++ {
			row := y - 1 + w%3
			col := i + k - 1
			m[row][col] = ch
			w++
		}
	}
	// draw borders and separators
	for i := 1; i <= 5; i++ {
		m[i-1][0] = 'A'
	}
	for i := 1; i <= 50; i++ {
		m[0][i-1] = 'A'
	}
	for i := 1; i <= 40; i++ {
		m[i-1][49] = 'A'
	}
	for i := 2; i <= 50; i++ {
		m[30][i-1] = 'A'
	}
	for i := 31; i <= 45; i++ {
		m[i-1][0] = 'D'
	}
	// fill extra segments
	if a > 90 {
		ff(3, a-90, 27, 'A')
		a = 90
	}
	if b > 90 {
		ff(14, b-90, 27, 'B')
		b = 90
	}
	if c > 90 {
		ff(25, c-90, 27, 'C')
		c = 90
	}
	if a > 45 {
		f(45, 7, 'A')
		a -= 45
	}
	f(a, 3, 'A')
	if b > 45 {
		f(45, 11, 'B')
		b -= 45
	}
	f(b, 15, 'B')
	if c > 45 {
		f(45, 19, 'C')
		c -= 45
	}
	f(c, 23, 'C')
	if d > 90 {
		ff(3, d-90, 32, 'D')
		d = 90
	}
	if d > 45 {
		ff(3, 45, 36, 'D')
		d -= 45
	}
	f(d, 45, 'D')
	// fill remaining
	for i := 1; i <= 31; i++ {
		for j := 1; j <= 50; j++ {
			if m[i-1][j-1] == '#' {
				m[i-1][j-1] = 'D'
			}
		}
	}
	for i := 31; i <= 50; i++ {
		for j := 1; j <= 50; j++ {
			if m[i-1][j-1] == '#' {
				m[i-1][j-1] = 'A'
			}
		}
	}
	// output
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, "50 50")
	for i := 0; i < 50; i++ {
		out.Write(m[i])
		out.WriteByte('\n')
	}
}
