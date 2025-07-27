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
		var px, py int
		fmt.Fscan(in, &px, &py)
		var s string
		fmt.Fscan(in, &s)
		var up, down, left, right int
		for i := 0; i < len(s); i++ {
			switch s[i] {
			case 'U':
				up++
			case 'D':
				down++
			case 'L':
				left++
			case 'R':
				right++
			}
		}
		ok := true
		if px > 0 {
			if right < px {
				ok = false
			}
		} else if px < 0 {
			if left < -px {
				ok = false
			}
		}
		if py > 0 {
			if up < py {
				ok = false
			}
		} else if py < 0 {
			if down < -py {
				ok = false
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
