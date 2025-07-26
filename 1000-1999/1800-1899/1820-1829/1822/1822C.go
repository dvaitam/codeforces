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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n int64
		fmt.Fscan(reader, &n)
		ans := n*n + (n-2)*(n-2)
		fmt.Fprintln(writer, ans)
	}
}
