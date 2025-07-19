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

	var t int
	fmt.Fscan(reader, &t)
	for i := 0; i < t; i++ {
		var n float64
		fmt.Fscan(reader, &n)
		ans := 1 / math.Tan(math.Pi/(2*n))
		fmt.Fprintf(writer, "%.6f\n", ans)
	}
}
