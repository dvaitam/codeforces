package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var T, n int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		fmt.Fscan(in, &n)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}
		k, kk := 0, 0
		for i := 1; i <= n; i++ {
			if a[i]%2 == 0 {
				k = 1
			} else {
				kk = 1
			}
		}
		if k*kk == 1 {
			fmt.Println(-1)
			continue
		}
		if kk == 1 {
			fmt.Println(30)
		} else {
			fmt.Println(31)
		}
		for i := 29; i >= 0; i-- {
			fmt.Print(1<<i, " ")
		}
		if k == 1 {
			fmt.Print(1, " ")
		}
		fmt.Println()
	}
}
