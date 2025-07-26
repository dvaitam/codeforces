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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)

		cnt := make([]int, 26)
		for _, ch := range s {
			cnt[ch-'a']++
		}
		pairs := 0
		for _, c := range cnt {
			pairs += c / 2
		}
		remaining := n - k
		if pairs >= remaining/2 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
