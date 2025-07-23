package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n uint64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if n%2 == 1 {
		fmt.Println(1)
	} else {
		fmt.Println(2)
	}
}
