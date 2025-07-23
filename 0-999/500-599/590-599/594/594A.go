package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Ints(a)
	half := n / 2
	ans := a[half] - a[0]
	for i := 1; i < half; i++ {
		d := a[i+half] - a[i]
		if d < ans {
			ans = d
		}
	}
	fmt.Println(ans)
}
