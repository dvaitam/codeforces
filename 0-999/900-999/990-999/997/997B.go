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

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	pre := []int64{0, 4, 10, 20, 35, 56, 83, 116, 155, 198, 244, 292}
	if n <= 11 {
		fmt.Fprintln(out, pre[n])
		return
	}
	fmt.Fprintln(out, 49*n-247)
}
