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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, a, b int
		fmt.Fscan(in, &n, &a, &b)
		var s string
		fmt.Fscan(in, &s)
		x, y := 0, 0
		found := false
		if x == a && y == b {
			found = true
		}
		for step := 0; step < 100 && !found; step++ {
			ch := s[step%len(s)]
			switch ch {
			case 'N':
				y++
			case 'E':
				x++
			case 'S':
				y--
			case 'W':
				x--
			}
			if x == a && y == b {
				found = true
				break
			}
		}
		if found {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
