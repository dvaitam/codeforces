package main

import (
	"bufio"
	"fmt"
	"os"
)

func build(ids []int, l, r int, s []byte, pos *int) []int {
	if r-l == 1 {
		return []int{ids[l]}
	}
	mid := (l + r) / 2
	left := build(ids, l, mid, s, pos)
	right := build(ids, mid, r, s, pos)
	i, j := 0, 0
	res := make([]int, 0, len(left)+len(right))
	for i < len(left) && j < len(right) {
		if *pos >= len(s) || s[*pos] == '0' {
			res = append(res, left[i])
			i++
		} else {
			res = append(res, right[j])
			j++
		}
		(*pos)++
	}
	res = append(res, left[i:]...)
	res = append(res, right[j:]...)
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	const n = 16
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = i
	}
	pos := 0
	order := build(ids, 0, n, []byte(s), &pos)
	ans := make([]int, n)
	for rank, idx := range order {
		ans[idx] = rank + 1
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, n)
	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
}
