package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxN = 200000

var prefix [maxN + 1]int64
var sumDigits [maxN + 1]int

func init() {
	for i := 1; i <= maxN; i++ {
		sumDigits[i] = sumDigits[i/10] + i%10
		prefix[i] = prefix[i-1] + int64(sumDigits[i])
	}
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
		if n > maxN {
			n = maxN
		}
		fmt.Fprintln(writer, prefix[n])
	}
}
