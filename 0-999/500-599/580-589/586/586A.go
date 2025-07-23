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
	ans := 0
	for i := 0; i < n; i++ {
		if a[i] == 1 {
			ans++
		} else if i > 0 && i < n-1 && a[i-1] == 1 && a[i+1] == 1 {
			ans++
		}
	}
	fmt.Println(ans)
}
