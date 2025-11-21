package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
		var n int64
		fmt.Fscan(reader, &n)

		s := strconv.FormatInt(n, 10)
		length := len(s)

		ans := int64(9 * (length - 1))
		base := int64(0)
		for i := 0; i < length; i++ {
			base = base*10 + 1
		}
		for digit := int64(1); digit <= 9; digit++ {
			if digit*base <= n {
				ans++
			}
		}

		fmt.Fprintln(writer, ans)
	}
}
