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
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			var c int
			fmt.Fscan(in, &c)
			freq[c]++
		}

		singles := 0
		multis := 0
		for _, f := range freq {
			if f == 1 {
				singles++
			} else {
				multis++
			}
		}

		moves := (n + 1) / 2
		takeSingles := (singles + 1) / 2
		if takeSingles > moves {
			takeSingles = moves
		}
		remaining := moves - takeSingles
		takeMultis := multis
		if takeMultis > remaining {
			takeMultis = remaining
		}

		ans := takeSingles*2 + takeMultis
		fmt.Fprintln(out, ans)
	}
}
