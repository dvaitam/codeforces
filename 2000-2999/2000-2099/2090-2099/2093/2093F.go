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
		var n, m int
		fmt.Fscan(in, &n, &m)

		a := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		covered := make([]bool, n)
		maxGood := 0

		for i := 0; i < m; i++ {
			good := 0
			for j := 0; j < n; j++ {
				var b string
				fmt.Fscan(in, &b)
				if b == a[j] {
					good++
					covered[j] = true
				}
			}
			if good > maxGood {
				maxGood = good
			}
		}

		ok := true
		for _, c := range covered {
			if !c {
				ok = false
				break
			}
		}

		if !ok {
			fmt.Fprintln(out, -1)
			continue
		}

		// Strategy: use a network with maxGood matches to fill all positions (n ops),
		// then fix the remaining n-maxGood positions one by one (clear + fill).
		ans := 3*n - 2*maxGood
		fmt.Fprintln(out, ans)
	}
}
