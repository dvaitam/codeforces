package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	var token string
	if _, err := fmt.Fscan(in, &token); err != nil {
		return
	}

	values := make([]int, 0, t)

	if token == "manual" {
		for i := 0; i < t; i++ {
			var w int
			fmt.Fscan(in, &w)
			values = append(values, w)
		}
	} else {
		// token is actually the first W
		first, _ := strconv.Atoi(token)
		values = append(values, first)
		for i := 1; i < t; i++ {
			var w int
			fmt.Fscan(in, &w)
			values = append(values, w)
		}
	}

	for i, w := range values {
		if i > 0 {
			fmt.Fprintln(out)
		}
		fmt.Fprint(out, w)
	}
}
