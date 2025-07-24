package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var x1, y1, x2, y2 int
	if _, err := fmt.Fscan(reader, &x1, &y1); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &x2, &y2); err != nil {
		return
	}

	dx := abs(x1 - x2)
	dy := abs(y1 - y2)

	ans := 2*(dx+dy) + 4
	if x1 == x2 || y1 == y2 {
		ans += 2
	}
	fmt.Fprintln(writer, ans)
}
