package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var a, b string
	fmt.Fscan(reader, &a)
	fmt.Fscan(reader, &b)
	if a == b {
		fmt.Println(-1)
	} else if len(a) > len(b) {
		fmt.Println(len(a))
	} else {
		fmt.Println(len(b))
	}
}
