package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(arr []int, x int) bool {
	i, j := 0, len(arr)-1
	for i < j {
		if arr[i] == arr[j] {
			i++
			j--
		} else if arr[i] == x {
			i++
		} else if arr[j] == x {
			j--
		} else {
			return false
		}
	}
	return true
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		l, r := 0, n-1
		for l < r && a[l] == a[r] {
			l++
			r--
		}
		if l >= r {
			fmt.Fprintln(out, "YES")
			continue
		}
		if check(a, a[l]) || check(a, a[r]) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
