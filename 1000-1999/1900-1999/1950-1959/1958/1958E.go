package main

import (
	"bufio"
	"fmt"
	"os"
)

func construct(n, k int) []int {
	if k == 1 {
		res := make([]int, n)
		for i := 0; i < n; i++ {
			res[i] = i + 1
		}
		return res
	}
	if k == 2 {
		if n < 3 {
			return nil
		}
		res := []int{n - 1, 1, n}
		used := map[int]bool{n - 1: true, 1: true, n: true}
		for i := 2; i <= n-2; i++ {
			if !used[i] {
				res = append(res, i)
			}
		}
		return res
	}
	if k == 3 {
		if n < 5 {
			return nil
		}
		res := []int{n - 1, 1, n - 2, 2, n}
		used := map[int]bool{n - 1: true, 1: true, n - 2: true, 2: true, n: true}
		for i := 3; i <= n; i++ {
			if !used[i] {
				res = append(res, i)
			}
		}
		return res
	}
	return nil
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		perm := construct(n, k)
		if perm == nil {
			fmt.Fprintln(out, -1)
		} else {
			for i, v := range perm {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, v)
			}
			fmt.Fprintln(out)
		}
	}
}
