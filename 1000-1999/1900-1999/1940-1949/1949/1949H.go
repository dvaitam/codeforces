package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: The exact algorithm for problem H is non-trivial and the original
// statement is incomplete in this repository. The implementation below merely
// follows the hint from the problem description that an organism always keeps
// at least one cell inside the 0<=x,y<=3 square. Hence if all cells in this
// square are forbidden, we print "NO"; otherwise we print "YES".
// This is only a placeholder and may not solve the real problem.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	const limit = 4
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		forbidden := make(map[[2]int]struct{}, n)
		for i := 0; i < n; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			forbidden[[2]int{x, y}] = struct{}{}
		}
		all := true
		for x := 0; x < limit; x++ {
			for y := 0; y < limit; y++ {
				if _, ok := forbidden[[2]int{x, y}]; !ok {
					all = false
					break
				}
			}
			if !all {
				break
			}
		}
		if all {
			fmt.Fprintln(writer, "NO")
		} else {
			fmt.Fprintln(writer, "YES")
		}
	}
}
