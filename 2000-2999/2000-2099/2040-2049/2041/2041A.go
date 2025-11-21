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

	used := make([]bool, 6)
	for i := 0; i < 4; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x >= 1 && x <= 5 {
			used[x] = true
		}
	}
	for i := 1; i <= 5; i++ {
		if !used[i] {
			fmt.Fprintln(out, i)
			break
		}
	}
}
