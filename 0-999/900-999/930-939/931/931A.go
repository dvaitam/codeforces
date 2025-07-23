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
	var a, b int
	if _, err := fmt.Fscan(in, &a, &b); err != nil {
		return
	}
	if a > b {
		a, b = b, a
	}
	diff := b - a
	x := diff / 2
	y := diff - x
	ans := x*(x+1)/2 + y*(y+1)/2
	fmt.Fprintln(out, ans)
}
