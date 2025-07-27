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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var a, k int64
		fmt.Fscan(reader, &a, &k)
		for i := int64(1); i < k; i++ {
			minD, maxD := int64(9), int64(0)
			temp := a
			for temp > 0 {
				d := temp % 10
				if d < minD {
					minD = d
				}
				if d > maxD {
					maxD = d
				}
				temp /= 10
			}
			if minD == 0 {
				break
			}
			a += minD * maxD
		}
		fmt.Fprintln(writer, a)
	}
}
