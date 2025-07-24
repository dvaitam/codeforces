package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	skills := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &skills[i])
	}

	arr := make([]pair, n)
	for i, v := range skills {
		arr[i] = pair{v, i}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].val == arr[j].val {
			return arr[i].idx < arr[j].idx
		}
		return arr[i].val < arr[j].val
	})

	ans := make([]int, n)
	countLess := 0
	for i := 0; i < n; i++ {
		if i > 0 && arr[i].val != arr[i-1].val {
			countLess = i
		}
		ans[arr[i].idx] = countLess
	}

	for i := 0; i < k; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		x--
		y--
		if skills[x] > skills[y] {
			ans[x]--
		} else if skills[y] > skills[x] {
			ans[y]--
		}
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, ans[i])
	}
	fmt.Fprintln(out)
}
