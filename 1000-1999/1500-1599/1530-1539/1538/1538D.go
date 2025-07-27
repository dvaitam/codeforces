package main

import (
	"bufio"
	"fmt"
	"os"
)

func primeCount(x int64) int {
	cnt := 0
	for x%2 == 0 {
		cnt++
		x /= 2
	}
	for i := int64(3); i*i <= x; i += 2 {
		for x%i == 0 {
			cnt++
			x /= i
		}
	}
	if x > 1 {
		cnt++
	}
	return cnt
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
		var a, b int64
		var k int
		fmt.Fscan(reader, &a, &b, &k)
		cntA := primeCount(a)
		cntB := primeCount(b)
		maxOps := cntA + cntB
		if k == 1 {
			if a != b && (a%b == 0 || b%a == 0) {
				fmt.Fprintln(writer, "Yes")
			} else {
				fmt.Fprintln(writer, "No")
			}
		} else {
			if k >= 2 && k <= maxOps {
				fmt.Fprintln(writer, "Yes")
			} else {
				fmt.Fprintln(writer, "No")
			}
		}
	}
}
