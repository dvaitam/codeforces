package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var n int
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	letters := []byte{'a', 'b', 'c', 'd'}
	for i := 0; i < n; i++ {
		writer.WriteByte(letters[i%4])
	}
	writer.WriteByte('\n')
}
