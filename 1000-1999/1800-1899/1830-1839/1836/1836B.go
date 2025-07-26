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
		var n, k, g int64
		fmt.Fscan(reader, &n, &k, &g)
		tVal := (g - 1) / 2
		maxSave := tVal * n
		total := k * g
		if maxSave > total {
			maxSave = total
		}
		ans := (maxSave / g) * g
		fmt.Fprintln(writer, ans)
	}
}
