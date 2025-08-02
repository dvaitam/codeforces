package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		if n <= 5 {
			fmt.Println(n/2 + 1)
			for i := 1; i <= n; i++ {
				fmt.Print(i/2+1, " ")
			}
			fmt.Println()
			continue
		}
		fmt.Println(4)
		for i := 1; i <= n; i++ {
			fmt.Print(i%4+1, " ")
		}
		fmt.Println()
	}
}
