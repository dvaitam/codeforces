package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a, b, c, d int
	if _, err := fmt.Fscan(in, &a); err != nil {
		return
	}
	fmt.Fscan(in, &b)
	fmt.Fscan(in, &c)
	fmt.Fscan(in, &d)
	fmt.Println(a ^ b ^ c ^ d)
}
