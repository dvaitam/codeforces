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
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		base := (n / k) * k
		rem := n % k
		half := k / 2
		if rem > half {
			rem = half
		}
		fmt.Fprintln(writer, base+rem)
	}
}
