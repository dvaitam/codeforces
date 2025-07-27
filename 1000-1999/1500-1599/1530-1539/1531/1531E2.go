package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	s   string
	pos int
	cur int
)

func countMemo(n int, memo map[int]int) int {
	if n <= 1 {
		return 0
	}
	if v, ok := memo[n]; ok {
		return v
	}
	l := n / 2
	r := n - l
	v := countMemo(l, memo) + countMemo(r, memo) + n - 1
	memo[n] = v
	return v
}

func restore(n int) []int {
	if n == 1 {
		cur++
		return []int{cur}
	}
	l := n / 2
	r := n - l
	left := restore(l)
	right := restore(r)
	i, j := 0, 0
	res := make([]int, 0, n)
	for i < len(left) && j < len(right) {
		if s[pos] == '0' {
			res = append(res, left[i])
			i++
		} else {
			res = append(res, right[j])
			j++
		}
		pos++
	}
	for i < len(left) {
		res = append(res, left[i])
		i++
	}
	for j < len(right) {
		res = append(res, right[j])
		j++
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fmt.Fscan(in, &s)
	L := len(s)
	memo := make(map[int]int)
	n := -1
	for i := 1; i <= 1000; i++ {
		if countMemo(i, memo) == L {
			n = i
			break
		}
	}
	if n == -1 {
		fmt.Println("-1")
		return
	}
	cur = 0
	pos = 0
	arr := restore(n)
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, arr[i])
	}
	out.WriteByte('\n')
	out.Flush()
}
