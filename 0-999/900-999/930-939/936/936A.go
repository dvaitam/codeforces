package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k, d, t int64
	if _, err := fmt.Fscan(in, &k, &d, &t); err != nil {
		return
	}

	cycleLen := ((k + d - 1) / d) * d
	cycleProgress := k + cycleLen
	totalNeed := 2 * t

	fullCycles := totalNeed / cycleProgress
	timeSpent := fullCycles * cycleLen
	remainder := totalNeed % cycleProgress

	extra := 0.0
	if remainder <= 2*k {
		extra = float64(remainder) / 2.0
	} else {
		extra = float64(k) + float64(remainder-2*k)
	}

	result := float64(timeSpent) + extra
	fmt.Printf("%.10f\n", result)
}
