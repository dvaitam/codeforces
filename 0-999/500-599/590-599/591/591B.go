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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	mapping := [26]byte{}
	for i := 0; i < 26; i++ {
		mapping[i] = byte('a' + i)
	}

	for i := 0; i < m; i++ {
		var xs, ys string
		fmt.Fscan(reader, &xs, &ys)
		x := xs[0]
		y := ys[0]
		for j := 0; j < 26; j++ {
			if mapping[j] == x {
				mapping[j] = y
			} else if mapping[j] == y {
				mapping[j] = x
			}
		}
	}

	bytes := []byte(s)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = mapping[bytes[i]-'a']
	}
	fmt.Fprintln(writer, string(bytes))
}
