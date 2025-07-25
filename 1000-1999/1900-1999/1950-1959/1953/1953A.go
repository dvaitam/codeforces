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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	for i := 0; i < n; i++ {
		var dummy float64
		fmt.Fscan(reader, &dummy)
	}

	b := strings.Builder{}
	b.WriteString("{d:")
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte(',')
		}
		b.WriteString(fmt.Sprintf("%d", i))
	}
	b.WriteByte('}')

	fmt.Fprint(writer, b.String())
}
