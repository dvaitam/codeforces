package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(w *bufio.Writer, r *bufio.Reader, i, j, k int) int {
	fmt.Fprintf(w, "? %d %d %d\n", i, j, k)
	w.Flush()
	var ans int
	fmt.Fscan(r, &ans)
	return ans
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		if n < 3 {
			fmt.Fprintln(writer, "! 1 1")
			writer.Flush()
			var verdict int
			fmt.Fscan(reader, &verdict)
			if verdict == -1 {
				return
			}
			continue
		}

		maxVal := -1
		pos := 3
		for i := 3; i <= n; i++ {
			res := query(writer, reader, 1, 2, i)
			if res > maxVal {
				maxVal = res
				pos = i
			}
		}

		maxVal = -1
		pos2 := 2
		for i := 2; i <= n; i++ {
			if i == pos {
				continue
			}
			res := query(writer, reader, 1, pos, i)
			if res > maxVal {
				maxVal = res
				pos2 = i
			}
		}

		fmt.Fprintf(writer, "! %d %d\n", pos, pos2)
		writer.Flush()
		var verdict int
		fmt.Fscan(reader, &verdict)
		if verdict == -1 {
			return
		}
	}
}
