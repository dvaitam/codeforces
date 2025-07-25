package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		counts := make([]int, 10)
		for i := 0; i < 10; i++ {
			fmt.Fscan(reader, &counts[i])
		}
		best := "1" + strings.Repeat("0", counts[0]+1)
		for d := 1; d <= 9; d++ {
			digit := fmt.Sprintf("%d", d)
			cand := strings.Repeat(digit, counts[d]+1)
			if len(cand) < len(best) || (len(cand) == len(best) && cand < best) {
				best = cand
			}
		}
		fmt.Fprintln(writer, best)
	}
}
