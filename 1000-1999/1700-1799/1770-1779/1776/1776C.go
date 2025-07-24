package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}
	sort.Slice(a, func(i, j int) bool { return a[i] > a[j] })

	l := m
	for _, d := range a {
		if d > l {
			fmt.Println("Bernardo")
			return
		}
		if d <= l/2 {
			l -= d
		} else {
			l /= 2
		}
	}
	fmt.Println("Alessia")
}
