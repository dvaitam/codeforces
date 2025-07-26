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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	ask := func(u, k int, s []int) int {
		fmt.Fprintf(writer, "? %d %d %d", u, k, len(s))
		for _, v := range s {
			fmt.Fprintf(writer, " %d", v)
		}
		fmt.Fprintln(writer)
		writer.Flush()
		var resp int
		if _, err := fmt.Fscan(reader, &resp); err != nil {
			os.Exit(0)
		}
		return resp
	}

	reachable := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if ask(i, n, []int{1}) == 1 {
			reachable = append(reachable, i)
		}
	}

	fmt.Fprintf(writer, "! %d", len(reachable))
	for _, v := range reachable {
		fmt.Fprintf(writer, " %d", v)
	}
	fmt.Fprintln(writer)
}
