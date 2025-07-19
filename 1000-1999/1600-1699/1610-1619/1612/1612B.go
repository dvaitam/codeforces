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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, a, b int
		fmt.Fscan(reader, &n, &a, &b)
		half := n / 2
		if a <= half+1 && b >= half {
			if a == half+1 && b == half {
				fmt.Fprintf(writer, "%d ", a)
				for i := half + 2; i <= n; i++ {
					fmt.Fprintf(writer, "%d ", i)
				}
				for i := 1; i < half; i++ {
					fmt.Fprintf(writer, "%d ", i)
				}
				fmt.Fprintf(writer, "%d\n", b)
			} else if a == half+1 || b == half {
				if b > a {
					fmt.Fprint(writer, "-1\n")
				}
			} else {
				fmt.Fprintf(writer, "%d ", a)
				for i := half + 1; i <= n; i++ {
					if i == a || i == b {
						continue
					}
					fmt.Fprintf(writer, "%d ", i)
				}
				for i := 1; i <= half; i++ {
					if i == a || i == b {
						continue
					}
					fmt.Fprintf(writer, "%d ", i)
				}
				fmt.Fprintf(writer, "%d\n", b)
			}
		} else {
			fmt.Fprint(writer, "-1\n")
		}
	}
}
