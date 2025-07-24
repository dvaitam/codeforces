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
	n := len(s)
	count := 0
	for i := 0; i < n; i++ {
		if s[i] == 'Q' {
			for j := i + 1; j < n; j++ {
				if s[j] == 'A' {
					for k := j + 1; k < n; k++ {
						if s[k] == 'Q' {
							count++
						}
					}
				}
			}
		}
	}
	fmt.Println(count)
}
