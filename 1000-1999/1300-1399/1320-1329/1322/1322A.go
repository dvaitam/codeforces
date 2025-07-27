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
	fmt.Fscan(reader, &n)
	var s string
	fmt.Fscan(reader, &s)

	if n%2 == 1 {
		fmt.Fprintln(writer, -1)
		return
	}

	count := 0
	for _, ch := range s {
		if ch == '(' {
			count++
		} else {
			count--
		}
	}
	if count != 0 {
		fmt.Fprintln(writer, -1)
		return
	}

	balance := 0
	start := -1
	ans := 0
	for i, ch := range s {
		if ch == '(' {
			balance++
		} else {
			balance--
		}
		if balance < 0 && start == -1 {
			start = i
		}
		if balance == 0 && start != -1 {
			ans += i - start + 1
			start = -1
		}
	}

	fmt.Fprintln(writer, ans)
}
