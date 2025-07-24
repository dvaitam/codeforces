package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var maxn, maxa, q int
	if _, err := fmt.Fscan(in, &maxn, &maxa, &q); err != nil {
		return
	}
	// TODO: implement full solution
	fmt.Println(0)
}
