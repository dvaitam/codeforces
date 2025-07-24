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
		var s string
		fmt.Fscan(reader, &s)
		sum := 0
		zeros := 0
		even := 0
		for i := 0; i < len(s); i++ {
			d := int(s[i] - '0')
			sum += d
			if d == 0 {
				zeros++
				even++
			} else if d%2 == 0 {
				even++
			}
		}
		if zeros > 0 && sum%3 == 0 && (zeros > 1 || even > 1) {
			fmt.Fprintln(writer, "red")
		} else {
			fmt.Fprintln(writer, "cyan")
		}
	}
}
