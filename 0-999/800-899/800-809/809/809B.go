package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	dishes := make([]int, 0, k)
	for i := 0; i < k; i++ {
		var x int
		if _, err := fmt.Fscan(in, &x); err != nil {
			break
		}
		dishes = append(dishes, x)
	}

	if len(dishes) >= 2 {
		fmt.Fprintf(out, "2 %d %d\n", dishes[0], dishes[1])
	} else if n >= 2 {
		fmt.Fprintf(out, "2 %d %d\n", 1, 2)
	} else {
		fmt.Fprintf(out, "2 %d %d\n", 1, 1)
	}
}
