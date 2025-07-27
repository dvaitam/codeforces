package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	sum := 0
	for _, r := range s {
		if unicode.IsDigit(r) {
			sum += int(r - '0')
		}
	}
	fmt.Println(sum)
}
