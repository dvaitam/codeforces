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
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		half := n / 2

		// build groups of equal scores
		groups := []int{}
		count := 1
		for i := 1; i < n; i++ {
			if arr[i] == arr[i-1] {
				count++
			} else {
				groups = append(groups, count)
				count = 1
			}
		}
		groups = append(groups, count)

		// take prefix of groups not exceeding half
		prefix := 0
		valid := []int{}
		for _, g := range groups {
			if prefix+g > half {
				break
			}
			valid = append(valid, g)
			prefix += g
		}

		if len(valid) < 3 {
			fmt.Fprintln(out, "0 0 0")
			continue
		}

		gold := valid[0]
		silver := 0
		i := 1
		for i < len(valid) && silver <= gold {
			silver += valid[i]
			i++
		}
		if silver <= gold || i >= len(valid) {
			fmt.Fprintln(out, "0 0 0")
			continue
		}

		bronze := prefix - gold - silver
		if bronze <= gold {
			fmt.Fprintln(out, "0 0 0")
			continue
		}

		fmt.Fprintf(out, "%d %d %d\n", gold, silver, bronze)
	}
}
