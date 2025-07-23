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
	f := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &f[i])
	}
	for i := 1; i <= n; i++ {
		j := f[i]
		if j < 1 || j > n {
			continue
		}
		k := f[j]
		if k < 1 || k > n {
			continue
		}
		if f[k] == i {
			fmt.Println("YES")
			return
		}
	}
	fmt.Println("NO")
}
