package main

import (
	"bufio"
	"fmt"
	"os"
)

const limit = 1000000

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	valid := make(map[int]struct{})
	a := 1
	step := 0
	for a <= limit {
		valid[a] = struct{}{}
		if a+1 <= limit {
			valid[a+1] = struct{}{}
		}
		if step%2 == 0 {
			a = 2*a + 2
		} else {
			a = 2*a + 1
		}
		step++
	}

	if _, ok := valid[n]; ok {
		fmt.Fprintln(writer, 1)
	} else {
		fmt.Fprintln(writer, 0)
	}
}
