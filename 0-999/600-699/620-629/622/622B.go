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

	var ts string
	fmt.Fscan(reader, &ts)

	var h, m, a int
	fmt.Sscanf(ts, "%d:%d", &h, &m)
	fmt.Fscan(reader, &a)

	total := h*60 + m + a
	total %= 24 * 60
	h = total / 60
	m = total % 60

	fmt.Fprintf(writer, "%02d:%02d\n", h, m)
}
