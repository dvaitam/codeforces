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

	var n, k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	half := n / 2
	diplomas := half / (k + 1)
	certificates := diplomas * k
	notWinners := n - diplomas - certificates

	fmt.Fprintf(writer, "%d %d %d\n", diplomas, certificates, notWinners)
}
