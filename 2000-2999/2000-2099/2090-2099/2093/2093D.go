package main

import (
	"bufio"
	"fmt"
	"os"
)

func idx(x, y int, n int) int64 {
	if n == 1 {
		switch {
		case x == 0 && y == 0:
			return 0
		case x == 1 && y == 1:
			return 1
		case x == 1 && y == 0:
			return 2
		default:
			return 3
		}
	}
	half := 1 << (n - 1)
	block := int64(1) << (2 * (n - 1))
	if x < half && y < half {
		return idx(x, y, n-1)
	} else if x >= half && y >= half {
		return block + idx(x-half, y-half, n-1)
	} else if x >= half && y < half {
		return 2*block + idx(x-half, y, n-1)
	}
	return 3*block + idx(x, y-half, n-1)
}

func coord(d int64, n int) (int, int) {
	if n == 1 {
		switch d {
		case 0:
			return 0, 0
		case 1:
			return 1, 1
		case 2:
			return 1, 0
		default:
			return 0, 1
		}
	}
	block := int64(1) << (2 * (n - 1))
	half := 1 << (n - 1)
	order := d / block
	rem := d % block
	x, y := coord(rem, n-1)
	switch order {
	case 0:
		return x, y
	case 1:
		return x + half, y + half
	case 2:
		return x + half, y
	default:
		return x, y + half
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var typ string
			fmt.Fscan(in, &typ)
			if typ == "->" {
				var x, y int
				fmt.Fscan(in, &x, &y)
				val := idx(x-1, y-1, n) + 1
				fmt.Fprintln(out, val)
			} else {
				var d int64
				fmt.Fscan(in, &d)
				x, y := coord(d-1, n)
				fmt.Fprintf(out, "%d %d\n", x+1, y+1)
			}
		}
	}
}
