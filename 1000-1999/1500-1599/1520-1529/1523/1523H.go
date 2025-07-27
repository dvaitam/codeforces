package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	pos, rem int
}

func solve(arr []int, k int) int {
	n := len(arr)
	if n == 1 {
		return 0
	}
	visited := make([][]bool, k+1)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	cur := []state{{0, 0}}
	visited[0][0] = true
	steps := 0
	for len(cur) > 0 {
		next := make([]state, 0)
		for _, st := range cur {
			if st.pos == n-1 {
				return steps
			}
			avail := k - st.rem
			limit := st.pos + arr[st.pos] + avail
			if limit >= n {
				limit = n - 1
			}
			for j := st.pos + 1; j <= limit; j++ {
				need := 0
				if j-st.pos > arr[st.pos] {
					need = j - st.pos - arr[st.pos]
				}
				nr := st.rem + need
				if nr <= k && !visited[nr][j] {
					visited[nr][j] = true
					next = append(next, state{j, nr})
				}
			}
		}
		cur = next
		steps++
	}
	return -1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	fmt.Fscan(reader, &n, &q)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for ; q > 0; q-- {
		var l, r, k int
		fmt.Fscan(reader, &l, &r, &k)
		sub := make([]int, r-l+1)
		copy(sub, a[l-1:r])
		ans := solve(sub, k)
		fmt.Fprintln(writer, ans)
	}
}
