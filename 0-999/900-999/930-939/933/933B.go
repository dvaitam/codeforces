package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var p, k int64
	if _, err := fmt.Fscan(reader, &p, &k); err != nil {
		return
	}

	digits := make([]int64, 0)
	for p != 0 {
		r := ((p % k) + k) % k
		digits = append(digits, r)
		p = (p - r) / -k
	}

	if len(digits) == 0 {
		digits = append(digits, 0)
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, len(digits))
	for i, v := range digits {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
	fmt.Fprintln(writer)
	writer.Flush()
}
