package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for ; T > 0; T-- {
		var x string
		var k int
		fmt.Fscan(in, &x)
		fmt.Fscan(in, &k)
		n := len(x)
		L := n - k
		positions := make([][]int, 10)
		for i, c := range x {
			d := int(c - '0')
			positions[d] = append(positions[d], i)
		}
		ptr := make([]int, 10)
		res := make([]byte, 0, L)
		start := 0
		for j := 0; j < L; j++ {
			end := n - (L - j)
			digitStart := 0
			if j == 0 {
				digitStart = 1
			}
			chosenDigit := -1
			chosenIndex := -1
			for d := digitStart; d <= 9; d++ {
				p := ptr[d]
				for p < len(positions[d]) && positions[d][p] < start {
					p++
				}
				ptr[d] = p
				if p < len(positions[d]) {
					idx := positions[d][p]
					if idx <= end {
						chosenDigit = d
						chosenIndex = idx
						break
					}
				}
			}
			res = append(res, byte(chosenDigit+'0'))
			start = chosenIndex + 1
		}
		fmt.Fprintln(writer, string(res))
	}
}
