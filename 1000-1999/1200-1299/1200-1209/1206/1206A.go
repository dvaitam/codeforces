package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(reader, &a[i])
	}

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	b := make([]int, m)
	for i := range b {
		fmt.Fscan(reader, &b[i])
	}

	ma := make(map[int]bool, len(a))
	mb := make(map[int]bool, len(b))
	for _, v := range a {
		ma[v] = true
	}
	for _, v := range b {
		mb[v] = true
	}

	for _, x := range a {
		for _, y := range b {
			s := x + y
			if !ma[s] && !mb[s] {
				fmt.Println(x, y)
				return
			}
		}
	}
}
