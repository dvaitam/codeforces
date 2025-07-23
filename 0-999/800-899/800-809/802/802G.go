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
	target := "heidi"
	j := 0
	for i := 0; i < len(s) && j < len(target); i++ {
		if s[i] == target[j] {
			j++
		}
	}
	if j == len(target) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
