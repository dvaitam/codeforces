package main

import (
	"bufio"
	"fmt"
	"os"
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
	if n < 3 {
		fmt.Print(0)
		return
	}
	cnt := 0
	for i := 1; i < n-1; i++ {
		if (a[i] > a[i-1] && a[i] > a[i+1]) || (a[i] < a[i-1] && a[i] < a[i+1]) {
			cnt++
		}
	}
	fmt.Print(cnt)
}
