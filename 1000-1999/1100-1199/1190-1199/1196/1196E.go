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
	var T int
	fmt.Fscan(reader, &T)
	for T > 0 {
		T--
		var b, w int
		fmt.Fscan(reader, &b, &w)
		det := 0
		if b < w {
			b, w = w, b
			det = 1
		}
		if b-(w-1) > w*2+2 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		// Place stones in a line and surround
		for i := 1; i <= w; i++ {
			if i != 1 {
				fmt.Fprintf(writer, "%d %d\n", (i*2-1)+det, 2)
			}
			fmt.Fprintf(writer, "%d %d\n", (i*2)+det, 2)
		}
		b -= (w - 1)
		if b > 0 {
			fmt.Fprintf(writer, "%d %d\n", 1+det, 2)
			b--
		}
		for i := 1; i <= w; i++ {
			if b == 0 {
				break
			}
			fmt.Fprintf(writer, "%d %d\n", (i*2)+det, 1)
			b--
			if b == 0 {
				break
			}
			fmt.Fprintf(writer, "%d %d\n", (i*2)+det, 3)
			b--
		}
		if b > 0 {
			fmt.Fprintf(writer, "%d %d\n", (1+w*2)+det, 2)
			// b-- // not needed after last
		}
	}
}
