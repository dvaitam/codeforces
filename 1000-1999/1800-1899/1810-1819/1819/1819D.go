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
		shops := make([][]int, n)
		unknown := make([]bool, n)
		for i := 0; i < n; i++ {
			var k int
			fmt.Fscan(in, &k)
			if k == 0 {
				unknown[i] = true
			} else {
				arr := make([]int, k)
				for j := 0; j < k; j++ {
					fmt.Fscan(in, &arr[j])
				}
				shops[i] = arr
			}
		}

		seen := make([]bool, m+1)
		unionCnt := 0
		unknownExist := false
		for i := n - 1; i >= 0; i-- {
			if unknown[i] {
				unknownExist = true
				continue
			}
			arr := shops[i]
			conflict := false
			for _, v := range arr {
				if seen[v] {
					conflict = true
					break
				}
			}
			if conflict {
				break
			}
			for _, v := range arr {
				if !seen[v] {
					seen[v] = true
					unionCnt++
				}
			}
		}

		if unknownExist {
			fmt.Fprintln(out, m)
		} else {
			fmt.Fprintln(out, unionCnt)
		}
	}
}
