package main

import (
	"bufio"
	"fmt"
	"os"
)

func possible(a []int) bool {
	st := []int{a[0]}
	for i := 1; i < len(a); i++ {
		v := a[i]
		for len(st) >= 2 && st[len(st)-1] > v && st[len(st)-1]-1 == st[len(st)-2] {
			st = st[:len(st)-1]
		}
		top := st[len(st)-1]
		if top < v {
			st = append(st, v)
		} else if top-1 == v {
			st[len(st)-1] = v
		} else if top == v {
			if v == 0 {
				return false
			}
		} else {
			return false
		}
	}
	for len(st) >= 2 && st[len(st)-1]-1 == st[len(st)-2] {
		st = st[:len(st)-1]
	}
	return len(st) == 1 && st[0] == 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		if possible(arr) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
