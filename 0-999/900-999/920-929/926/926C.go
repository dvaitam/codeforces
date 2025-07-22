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

	runs := []int{}
	var prev, val int
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &val)
		if i == 0 {
			prev = val
			runs = append(runs, 1)
		} else {
			if val == prev {
				runs[len(runs)-1]++
			} else {
				prev = val
				runs = append(runs, 1)
			}
		}
	}

	ok := true
	for i := 1; i < len(runs); i++ {
		if runs[i] != runs[0] {
			ok = false
			break
		}
	}

	if ok {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
