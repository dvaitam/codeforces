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

	gifts := 0
	run := 0
	for i := 0; i < n; i++ {
		var grade int
		fmt.Fscan(reader, &grade)
		if grade >= 4 {
			run++
		} else {
			gifts += run / 3
			run = 0
		}
	}
	gifts += run / 3
	fmt.Fprintln(writer, gifts)
}
