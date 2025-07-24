package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var res int64
	if n%2 == 0 {
		res = n / 2
	} else {
		res = -((n + 1) / 2)
	}
	fmt.Print(res)
}
