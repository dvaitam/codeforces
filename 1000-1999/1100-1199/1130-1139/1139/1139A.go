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
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	var res int64
	for i := 0; i < n && i < len(s); i++ {
		if (s[i]-'0')%2 == 0 {
			res += int64(i + 1)
		}
	}
	fmt.Fprint(writer, res)
}
