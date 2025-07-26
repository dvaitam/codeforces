package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)

		if t == "a" {
			fmt.Fprintln(writer, 1)
			continue
		}
		if strings.Contains(t, "a") {
			fmt.Fprintln(writer, -1)
			continue
		}

		ans := int64(1)
		for i := 0; i < len(s); i++ {
			ans *= 2
		}
		fmt.Fprintln(writer, ans)
	}
}
