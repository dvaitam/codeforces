package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	zeros := 0
	ones := 0
	for _, c := range s {
		if c == '0' {
			zeros++
		} else if c == '1' {
			ones++
		}
	}
	diff := zeros - ones
	if diff < 0 {
		diff = -diff
	}
	fmt.Println(diff)
}
