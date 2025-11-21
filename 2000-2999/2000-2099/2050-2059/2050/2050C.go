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
		var s string
		fmt.Fscan(in, &s)
		sumMod := 0
		count2 := 0
		count3 := 0
		for i := 0; i < len(s); i++ {
			d := int(s[i] - '0')
			sumMod = (sumMod + d) % 9
			if d == 2 {
				count2++
			} else if d == 3 {
				count3++
			}
		}

		target := (9 - sumMod%9) % 9
		if target == 0 {
			fmt.Fprintln(out, "YES")
			continue
		}

		eff2 := count2
		if eff2 > 9 {
			eff2 = 9
		}
		eff3 := count3
		if eff3 > 9 {
			eff3 = 9
		}

		possible := make([]bool, 9)
		possible[0] = true
		for i := 0; i < eff2; i++ {
			next := make([]bool, 9)
			copy(next, possible)
			for m := 0; m < 9; m++ {
				if possible[m] {
					next[(m+2)%9] = true
				}
			}
			possible = next
		}
		for i := 0; i < eff3; i++ {
			next := make([]bool, 9)
			copy(next, possible)
			for m := 0; m < 9; m++ {
				if possible[m] {
					next[(m+6)%9] = true
				}
			}
			possible = next
		}

		if possible[target] {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
