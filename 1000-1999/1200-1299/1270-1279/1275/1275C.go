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

	var n int
	fmt.Fscan(reader, &n)
	digits := []int64{0, 1, 3, 4}
	for i := 0; i < n; i++ {
		var pos int64
		fmt.Fscan(reader, &pos)
		var sum int64
		for j := 0; j < 25; j++ {
			sum += digits[pos%4]
			pos /= 4
		}
		fmt.Fprintln(writer, sum)
	}
}
