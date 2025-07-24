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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	children := make([]int, n)
	for i := 0; i < n; i++ {
		children[i] = i + 1
	}
	leader := 0
	for i := 0; i < k; i++ {
		var a int
		fmt.Fscan(reader, &a)
		idx := (leader + a) % len(children)
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, children[idx])
		children = append(children[:idx], children[idx+1:]...)
		if len(children) > 0 {
			leader = idx % len(children)
		} else {
			leader = 0
		}
	}
	fmt.Fprintln(writer)
}
