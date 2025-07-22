package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	count := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			if (c-'0')%2 == 1 {
				count++
			}
		} else if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
			count++
		}
	}
	fmt.Println(count)
}
