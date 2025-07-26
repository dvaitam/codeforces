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
		var n int
		fmt.Fscan(reader, &n)
		bigrams := make([]string, n-2)
		for i := 0; i < n-2; i++ {
			fmt.Fscan(reader, &bigrams[i])
		}

		res := bigrams[0]
		for i := 1; i < n-2; i++ {
			if bigrams[i-1][1] != bigrams[i][0] {
				res += string(bigrams[i][0])
			}
			res += string(bigrams[i][1])
		}
		if len(res) < n {
			res += "a"
		}
		fmt.Fprintln(writer, res)
	}
}
