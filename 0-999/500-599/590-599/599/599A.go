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

	var d1, d2, d3 int
	if _, err := fmt.Fscan(reader, &d1, &d2, &d3); err != nil {
		return
	}

	option1 := d1 + d2 + d3
	option2 := 2 * (d1 + d2)
	option3 := 2 * (d1 + d3)
	option4 := 2 * (d2 + d3)

	ans := option1
	if option2 < ans {
		ans = option2
	}
	if option3 < ans {
		ans = option3
	}
	if option4 < ans {
		ans = option4
	}

	fmt.Fprintln(writer, ans)
}
