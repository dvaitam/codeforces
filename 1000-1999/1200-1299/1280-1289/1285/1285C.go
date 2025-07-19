package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var x int64
	fmt.Fscan(reader, &x)
	orig := x
	var v []int64
	for i := int64(2); i <= x/i; i++ {
		if x%i == 0 {
			var js int64 = 1
			for x%i == 0 {
				x /= i
				js *= i
			}
			v = append(v, js)
		}
	}
	if x > 1 {
		v = append(v, x)
	}

	// initial answer
	ansA, ansB := int64(1), orig
	n := len(v)
	// iterate subsets
	total := 1 << n
	for mask := 1; mask < total; mask++ {
		var a, b int64 = 1, 1
		for j := 0; j < n; j++ {
			if mask&(1<<j) != 0 {
				a *= v[j]
			} else {
				b *= v[j]
			}
		}
		var curMax, ansMax int64
		if a > b {
			curMax = a
		} else {
			curMax = b
		}
		if ansA > ansB {
			ansMax = ansA
		} else {
			ansMax = ansB
		}
		if curMax < ansMax {
			ansA, ansB = a, b
		}
	}

	// output in sorted order
	if ansA < ansB {
		fmt.Fprintln(writer, ansA, ansB)
	} else {
		fmt.Fprintln(writer, ansB, ansA)
	}
}
