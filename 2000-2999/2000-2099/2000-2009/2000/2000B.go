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
		var n int
		fmt.Fscan(in, &n)
		seats := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &seats[i])
			seats[i]--
		}
		occupied := make([]bool, n)
		ok := true
		for i := 0; i < n; i++ {
			pos := seats[i]
			if pos < 0 || pos >= n {
				ok = false
				continue
			}
			if i == 0 {
				occupied[pos] = true
				continue
			}
			valid := false
			if pos-1 >= 0 && occupied[pos-1] {
				valid = true
			}
			if pos+1 < n && occupied[pos+1] {
				valid = true
			}
			if !valid {
				ok = false
			}
			occupied[pos] = true
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
