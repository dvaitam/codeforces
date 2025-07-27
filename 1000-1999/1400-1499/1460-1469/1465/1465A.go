package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		cnt := 0
		for i := len(s) - 1; i >= 0 && s[i] == ')'; i-- {
			cnt++
		}
		if cnt > n-cnt {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
	}
}
