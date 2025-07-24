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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int64
		fmt.Fscan(reader, &n)
		a, b := int64(0), int64(0)
		i := int64(1)
		for n > 0 {
			take := i
			if take > n {
				take = n
			}
			if i%4 == 1 || i%4 == 0 {
				a += take
			} else {
				b += take
			}
			n -= take
			i++
		}
		fmt.Fprintf(writer, "%d %d\n", a, b)
	}
}
