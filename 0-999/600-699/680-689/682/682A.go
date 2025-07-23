package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int64
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	cntN := [5]int64{}
	cntM := [5]int64{}

	cntN[0] = n / 5
	cntM[0] = m / 5
	for i := int64(1); i < 5; i++ {
		cntN[i] = n / 5
		if n%5 >= i {
			cntN[i]++
		}
		cntM[i] = m / 5
		if m%5 >= i {
			cntM[i]++
		}
	}
	var ans int64
	for r := 0; r < 5; r++ {
		comp := (5 - r) % 5
		ans += cntN[r] * cntM[comp]
	}
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, ans)
}
