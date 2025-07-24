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
	var k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	for i := 0; i < k; i++ {
		if n%10 == 0 {
			n /= 10
		} else {
			n--
		}
	}
	fmt.Fprintln(out, n)
}
