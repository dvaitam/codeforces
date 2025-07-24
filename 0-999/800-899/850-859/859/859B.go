package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	best := int(1<<31 - 1)
	limit := int(math.Sqrt(float64(n)))
	for w := 1; w <= limit; w++ {
		h := (n + w - 1) / w
		p := 2 * (w + h)
		if p < best {
			best = p
		}
	}
	fmt.Println(best)
}
