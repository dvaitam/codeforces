package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemC.txt for contest 1851.
// It determines whether one can form a path from the first tile to the last
// using jumps to the right so that the path length is divisible by k and every
// block of k visited tiles has the same colour. We only need to know the
// positions of the k-th occurrence of the first tile's colour and the k-th
// occurrence of the last tile's colour.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}

		if k == 1 {
			fmt.Fprintln(out, "YES")
			continue
		}

		first := c[0]
		last := c[n-1]

		if first == last {
			cnt := 0
			for _, v := range c {
				if v == first {
					cnt++
				}
			}
			if cnt >= k {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
			continue
		}

		pos1 := -1
		cnt := 0
		for i, v := range c {
			if v == first {
				cnt++
				if cnt == k {
					pos1 = i
					break
				}
			}
		}
		if pos1 == -1 {
			fmt.Fprintln(out, "NO")
			continue
		}

		pos2 := -1
		cnt = 0
		for i := n - 1; i >= 0; i-- {
			if c[i] == last {
				cnt++
				if cnt == k {
					pos2 = i
					break
				}
			}
		}
		if pos2 == -1 {
			fmt.Fprintln(out, "NO")
			continue
		}

		if pos1 < pos2 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
