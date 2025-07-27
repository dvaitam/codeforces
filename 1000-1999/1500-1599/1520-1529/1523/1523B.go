package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(reader, &x)
		}

		ops := make([][3]int, 0, 3*n)
		for i := 1; i <= n; i += 2 {
			j := i + 1
			ops = append(ops, [3]int{1, i, j})
			ops = append(ops, [3]int{2, i, j})
			ops = append(ops, [3]int{1, i, j})
			ops = append(ops, [3]int{1, i, j})
			ops = append(ops, [3]int{2, i, j})
			ops = append(ops, [3]int{1, i, j})
		}

		fmt.Fprintln(writer, len(ops))
		for _, op := range ops {
			fmt.Fprintf(writer, "%d %d %d\n", op[0], op[1], op[2])
		}
	}
}
