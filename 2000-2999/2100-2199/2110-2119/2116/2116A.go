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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b, c, d int64
		fmt.Fscan(reader, &a, &b, &c, &d)
		var res string
		if d <= c {
			res = "Gellyfish"
		} else if c < d && a <= (d-c) {
			res = "Flower"
		} else {
			res = "Gellyfish"
		}
		fmt.Fprintln(writer, res)
	}
}
