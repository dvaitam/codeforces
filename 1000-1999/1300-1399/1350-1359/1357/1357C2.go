package main

import (
	"bufio"
	"fmt"
	"math"
	"math/bits"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, parity int
	if _, err := fmt.Fscan(reader, &n, &parity); err != nil {
		return
	}
	if parity != 0 {
		parity = 1
	}

	count := 0
	for mask := 0; mask < 1<<uint(n); mask++ {
		if bits.OnesCount(uint(mask))%2 == parity {
			count++
		}
	}
	amplitude := 1.0 / math.Sqrt(float64(count))
	fmt.Fprintf(writer, "Amplitude per state: %.6f\n", amplitude)
	for mask := 0; mask < 1<<uint(n); mask++ {
		if bits.OnesCount(uint(mask))%2 == parity {
			fmt.Fprintf(writer, "%0*b\n", n, mask)
		}
	}
}
