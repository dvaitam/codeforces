package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var p, y int
	if _, err := fmt.Fscan(reader, &p, &y); err != nil {
		return
	}
	for cand := y; cand > p; cand-- {
		limit := int(math.Sqrt(float64(cand)))
		ok := true
		for d := 2; d <= p && d <= limit; d++ {
			if cand%d == 0 {
				ok = false
				break
			}
		}
		if ok {
			fmt.Println(cand)
			return
		}
	}
	fmt.Println(-1)
}
