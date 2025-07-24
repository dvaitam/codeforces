package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	fmt.Fscan(in, &n, &k)

	s := n - 1
	b := s / k
	r := s % k
	d1 := 2 * b
	if r == 1 {
		d1 += 1
	} else if r >= 2 {
		d1 += 2
	}
	d2 := n - k + 1

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if d1 < d2 {
		fmt.Fprintln(out, d1)
		node := 2
		for i := 0; i < k; i++ {
			length := b
			if i < r {
				length++
			}
			cur := 1
			for j := 0; j < length; j++ {
				fmt.Fprintln(out, cur, node)
				cur = node
				node++
			}
		}
	} else {
		fmt.Fprintln(out, d2)
		L := n - k + 1
		for i := 1; i <= L; i++ {
			fmt.Fprintln(out, i, i+1)
		}
		mid := (L + 2) / 2
		node := L + 2
		for i := 0; i < k-2; i++ {
			fmt.Fprintln(out, mid, node)
			node++
		}
	}
}
