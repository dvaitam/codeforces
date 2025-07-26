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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		members := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &members[i])
		}
		base := members[0]
		count := 0
		for i := 0; i < n; i++ {
			if members[i] == base {
				count++
			}
		}
		fmt.Fprintln(writer, count)
	}
}
