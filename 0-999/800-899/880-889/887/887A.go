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
	idx := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			idx = i
			break
		}
	}
	if idx == -1 {
		fmt.Print("no")
		return
	}
	zero := 0
	for i := idx + 1; i < len(s); i++ {
		if s[i] == '0' {
			zero++
		}
	}
	if zero >= 6 {
		fmt.Print("yes")
	} else {
		fmt.Print("no")
	}
}
