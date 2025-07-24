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
	var home string
	fmt.Fscan(reader, &home)

	diff := 0
	for i := 0; i < n; i++ {
		var flight string
		fmt.Fscan(reader, &flight)
		if len(flight) < 7 {
			continue
		}
		if flight[:3] == home {
			diff++
		} else {
			diff--
		}
	}

	if diff == 0 {
		fmt.Fprintln(writer, "home")
	} else {
		fmt.Fprintln(writer, "contest")
	}
}
