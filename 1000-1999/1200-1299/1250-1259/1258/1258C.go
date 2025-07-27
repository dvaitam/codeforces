package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b int64
	_, err := fmt.Fscan(reader, &a, &b)
	if err != nil {
		// error reading input; exit
		return
	}
	if a > b {
		fmt.Println(a)
	} else {
		fmt.Println(b)
	}
}
