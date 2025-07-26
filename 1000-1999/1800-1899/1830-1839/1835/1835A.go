package main

import (
	"bufio"
	"fmt"
	"os"
)

func pow10(n int) int64 {
	res := int64(1)
	for i := 0; i < n; i++ {
		res *= 10
	}
	return res
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
		var A, B, C int
		var k int64
		fmt.Fscan(reader, &A, &B, &C, &k)

		aStart := pow10(A - 1)
		aEnd := pow10(A) - 1
		bMin := pow10(B - 1)
		bMax := pow10(B) - 1
		cMin := pow10(C - 1)
		cMax := pow10(C) - 1

		found := false
		var ansA, ansB int64
		for a := aStart; a <= aEnd; a++ {
			bStart := bMin
			if cMin-a > bStart {
				bStart = cMin - a
			}
			bEnd := bMax
			if cMax-a < bEnd {
				bEnd = cMax - a
			}
			if bStart <= bEnd {
				cnt := bEnd - bStart + 1
				if k > cnt {
					k -= cnt
				} else {
					ansA = a
					ansB = bStart + k - 1
					found = true
					break
				}
			}
		}
		if !found {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintf(writer, "%d + %d = %d\n", ansA, ansB, ansA+ansB)
		}
	}
}
