package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	var upperSum, lowerSum int
	for _, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			upperSum += int(ch - 'A' + 1)
		} else if ch >= 'a' && ch <= 'z' {
			lowerSum += int(ch - 'a' + 1)
		}
	}
	fmt.Println(upperSum - lowerSum)
}
