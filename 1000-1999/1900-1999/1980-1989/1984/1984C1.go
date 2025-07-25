package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		posMin, posMax := int64(0), int64(0)
		hasPos := true
		negMin, negMax := int64(0), int64(0)
		hasNeg := true

		for _, v := range a {
			newPos := make([]int64, 0, 8)
			newNeg := make([]int64, 0, 8)
			if hasPos {
				vals := []int64{posMin}
				if posMax != posMin {
					vals = append(vals, posMax)
				}
				for _, p := range vals {
					x := p + v
					if x >= 0 {
						newPos = append(newPos, x)
					} else {
						newNeg = append(newNeg, x)
					}
					newPos = append(newPos, abs64(x))
				}
			}
			if hasNeg {
				vals := []int64{negMin}
				if negMax != negMin {
					vals = append(vals, negMax)
				}
				for _, nval := range vals {
					x := nval + v
					if x >= 0 {
						newPos = append(newPos, x)
					} else {
						newNeg = append(newNeg, x)
					}
					newPos = append(newPos, abs64(x))
				}
			}

			if len(newPos) > 0 {
				posMin, posMax = newPos[0], newPos[0]
				for _, val := range newPos[1:] {
					if val < posMin {
						posMin = val
					}
					if val > posMax {
						posMax = val
					}
				}
				hasPos = true
			} else {
				hasPos = false
			}

			if len(newNeg) > 0 {
				negMin, negMax = newNeg[0], newNeg[0]
				for _, val := range newNeg[1:] {
					if val < negMin {
						negMin = val
					}
					if val > negMax {
						negMax = val
					}
				}
				hasNeg = true
			} else {
				hasNeg = false
			}
		}

		result := int64(0)
		hasResult := false
		if hasPos {
			result = posMax
			hasResult = true
		}
		if hasNeg {
			if !hasResult || abs64(negMin) > result {
				result = abs64(negMin)
				hasResult = true
			}
			if abs64(negMax) > result {
				result = abs64(negMax)
				hasResult = true
			}
		}
		if !hasResult {
			result = 0
		}
		fmt.Fprintln(writer, result)
	}
}
