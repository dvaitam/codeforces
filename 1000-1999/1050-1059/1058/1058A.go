package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	hard := false
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x == 1 {
			hard = true
		}
	}
	if hard {
		fmt.Println("HARD")
	} else {
		fmt.Println("EASY")
	}
}
