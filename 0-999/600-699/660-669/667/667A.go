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

	var d, h, v, e float64
	if _, err := fmt.Fscan(reader, &d, &h, &v, &e); err != nil {
		return
	}

	area := math.Pi * d * d / 4.0
	// If rain fills the cup faster or equal to our drinking speed, it's impossible
	if v <= area*e {
		fmt.Fprintln(writer, "NO")
		return
	}
	// Compute time until the cup is empty
	t := h / (v/area - e)
	fmt.Fprintln(writer, "YES")
	fmt.Fprintf(writer, "%.12f\n", t)
}
