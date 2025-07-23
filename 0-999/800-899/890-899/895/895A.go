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
	angles := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &angles[i])
	}

	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + angles[i]
	}

	best := 360
	for i := 0; i < n; i++ {
		for j := i; j <= n; j++ {
			s := prefix[j] - prefix[i]
			diff := int(math.Abs(float64(360 - 2*s)))
			if diff < best {
				best = diff
			}
		}
	}
	fmt.Println(best)
}
