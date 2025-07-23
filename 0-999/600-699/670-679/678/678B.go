package main

import (
	"bufio"
	"fmt"
	"os"
)

func isLeap(y int) bool {
	if y%400 == 0 {
		return true
	}
	if y%100 == 0 {
		return false
	}
	return y%4 == 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var y int
	if _, err := fmt.Fscan(reader, &y); err != nil {
		return
	}

	origLeap := isLeap(y)
	shift := 0
	for {
		if isLeap(y) {
			shift = (shift + 2) % 7
		} else {
			shift = (shift + 1) % 7
		}
		y++
		if shift == 0 && isLeap(y) == origLeap {
			fmt.Fprintln(writer, y)
			return
		}
	}
}
