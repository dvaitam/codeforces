package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x int
		fmt.Fscan(in, &x)
		digitPresent := [10]bool{}
		temp := x
		for temp > 0 {
			digitPresent[temp%10] = true
			temp /= 10
		}
		ans := -1
		for y := 0; y <= 9; y++ {
			if digitPresent[y] {
				ans = y
				break
			}
		}
		if ans == -1 {
			fmt.Fprintln(out, 10)
		} else {
			fmt.Fprintln(out, ans)
		}
	}
}
