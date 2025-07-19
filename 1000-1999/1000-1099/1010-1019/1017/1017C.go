package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)

	m := int(math.Sqrt(float64(n)))
	top := n / m

	for i := m*top + 1; i <= n; i++ {
		fmt.Fprintf(writer, "%d ", i)
	}

	for s := m*top - m; s >= 0; s -= m {
		for i := 1; i <= m; i++ {
			fmt.Fprintf(writer, "%d ", s+i)
		}
	}

	fmt.Fprintln(writer)
}
