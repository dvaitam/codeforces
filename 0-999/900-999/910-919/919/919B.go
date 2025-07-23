package main

import (
	"bufio"
	"fmt"
	"os"
)

func sumDigits(x int) int {
	s := 0
	for x > 0 {
		s += x % 10
		x /= 10
	}
	return s
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	cnt := 0
	for i := 19; ; i++ {
		if sumDigits(i) == 10 {
			cnt++
			if cnt == k {
				fmt.Fprintln(out, i)
				return
			}
		}
	}
}
