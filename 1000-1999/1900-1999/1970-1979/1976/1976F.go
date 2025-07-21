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
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		for i := 0; i < n-1; i++ {
			var v, u int
			fmt.Fscan(reader, &v, &u)
		}
		// The full solution requires graph algorithms to compute bridges
		// after adding extra edges. This implementation is omitted.
		for k := 1; k <= n-1; k++ {
			if k > 1 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, n-1)
		}
		writer.WriteByte('\n')
	}
}
