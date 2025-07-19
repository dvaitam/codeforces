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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	primes := make([]int, 0, 16)
	for i := 2; i <= n; i++ {
		if i > 2 {
			writer.WriteByte(' ')
		}
		limit := int(math.Sqrt(float64(i)))
		first := 0
		found := false
		for j, p := range primes {
			if p > limit {
				break
			}
			if i%p == 0 {
				found = true
				first = j + 1
				break
			}
		}
		if !found {
			primes = append(primes, i)
			first = len(primes)
		}
		fmt.Fprint(writer, first)
	}
	writer.WriteByte('\n')
}
