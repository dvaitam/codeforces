package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k uint64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	l := uint64(1)
	limit := n + 1
	for i := uint64(2); i <= k && l <= limit; i++ {
		g := gcd(l, i)
		if l/g > limit/i {
			l = limit + 1
			break
		}
		l = l / g * i
	}

	if l <= limit && limit%l == 0 {
		fmt.Fprintln(writer, "Yes")
	} else {
		fmt.Fprintln(writer, "No")
	}
}
