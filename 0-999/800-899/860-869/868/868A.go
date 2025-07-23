package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var pass string
	if _, err := fmt.Fscan(reader, &pass); err != nil {
		return
	}
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	haveFirst := false
	haveSecond := false
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(reader, &w)
		if w == pass {
			fmt.Println("YES")
			return
		}
		if w[1] == pass[0] {
			haveFirst = true
		}
		if w[0] == pass[1] {
			haveSecond = true
		}
	}
	if haveFirst && haveSecond {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
