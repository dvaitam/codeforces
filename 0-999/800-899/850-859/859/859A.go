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

	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}
	maxRank := 0
	for i := 0; i < k; i++ {
		var r int
		fmt.Fscan(reader, &r)
		if r > maxRank {
			maxRank = r
		}
	}
	if maxRank < 25 {
		fmt.Fprintln(writer, 0)
	} else {
		fmt.Fprintln(writer, maxRank-25)
	}
}
