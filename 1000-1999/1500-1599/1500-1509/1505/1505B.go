package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var a int
	if _, err := fmt.Fscan(in, &a); err != nil {
		return
	}
	res := a % 9
	if res == 0 {
		res = 9
	}
	fmt.Fprintln(out, res)
}
