package main

import (
	"bufio"
	"fmt"
	"os"
)

func nextGood(n int64) int64 {
	for {
		t := n
		good := true
		for t > 0 {
			if t%3 == 2 {
				good = false
				break
			}
			t /= 3
		}
		if good {
			return n
		}
		n++
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int64
		fmt.Fscan(reader, &n)
		fmt.Fprintln(writer, nextGood(n))
	}
}
